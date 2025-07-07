package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	rpc "github.com/hihibug/micro_module/Framework/rpc/config"
	"github.com/hihibug/micro_module/core/etcd"
	etcdClientV3 "go.etcd.io/etcd/client/v3"
)

type ServiceLeaseKeepAliveResponse = *etcdClientV3.LeaseKeepAliveResponse

type Etcd struct {
	Srv               etcd.Etcd
	Conf              Config
	Ctx               context.Context
	CtxCancel         context.CancelFunc
	leaseID           etcdClientV3.LeaseID //租约ID
	SubscribeCallback func([]rpc.Config)
}

func NewEtcd(conf Config) *Etcd {
	cli, err := etcd.NewEtcd(&conf.Etcd)
	if err != nil {
		panic("registry etcd create client err :" + err.Error())
	}
	return &Etcd{
		Srv:  cli,
		Conf: conf,
	}
}

func (e *Etcd) Client() any {
	return e
}

func (e *Etcd) RegisterService(conf *rpc.Config) error {
	//设置租约时间
	resp, err := e.Srv.LeaseGrant(e.Conf.Etcd.TimeOut)
	if err != nil {
		return err
	}
	e.Conf.RegisterData = conf
	//注册服务并绑定租约
	key := "/" + conf.ServiceName
	if e.Conf.RegisterData.GroupName != "" {
		key = "/" + e.Conf.RegisterData.GroupName + key
		if e.Conf.RegisterData.ClusterName != "" {
			key = "/" + e.Conf.RegisterData.ClusterName + key
		}
	}

	val, _ := json.Marshal(e.Conf.RegisterData)
	_, err = e.Srv.PutLease(key+"/"+strconv.Itoa(int(resp.ID)), string(val), etcdClientV3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	//设置续租 定期发送需求请求
	leaseRespChan, err := e.Srv.Clients().KeepAlive(context.Background(), resp.ID)
	go func() {
		for {
			_ = <-leaseRespChan
			if e.SubscribeCallback != nil {
				prefix, err := e.Srv.GetPrefix(key)
				if err != nil {
					log.Println("SubscribeCallback Get Data Err :" + err.Error())
				}
				result := make([]rpc.Config, 0)
				for _, kv := range prefix.Kvs {
					val := rpc.Config{}
					_ = json.Unmarshal(kv.Value, &val)
					result = append(result, val)
				}
				e.SubscribeCallback(result)
			}
		}
	}()

	if err != nil {
		return err
	}
	e.leaseID = resp.ID

	log.Printf("rpc register: %s", "etcd")
	return nil
}

func (e *Etcd) DiscoveryService(sName, GroupName, ClusterName string) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	key := "/" + sName
	if GroupName != "" {
		key = "/" + GroupName + key
		if ClusterName != "" {
			key = "/" + ClusterName + key
		}
	}
	ipPool, err := IpPool(key, e.Srv.Clients(), ctx, cancel)
	if err != nil {
		return "", err
	}
	defer cancel()
	go ipPool.watcher()

	if ipPool.Height() == 0 {
		return "", errors.New(key + " does not exist")
	}

	return ipPool.WeightedRoundRobin(), nil
}

func (e *Etcd) DeregisterService() error {
	_, err := e.Srv.Clients().Revoke(context.Background(), e.leaseID)
	if err != nil {
		return err
	}
	return e.Srv.Close()
}

func (e *Etcd) SubscribeService(fn func([]rpc.Config)) {
	if e.Conf.RegisterData == nil {
		log.Printf("service not register")
		return
	}
	e.SubscribeCallback = fn
}

func (e *Etcd) UnsubscribeService() {
	e.SubscribeCallback = nil
}
