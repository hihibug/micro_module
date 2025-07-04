package etcd

import (
	"context"
	"encoding/json"
	rpcConf "github.com/hihibug/microdule/v2/Framework/rpc/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strconv"
	"sync"
)

type Pool struct {
	Ctx         context.Context
	CtxCancel   context.CancelFunc
	ServiceName string
	IpPool      *sync.Map
	EtcdClient  *clientv3.Client
}

func IpPool(sName string, cli *clientv3.Client, ctx context.Context, cancel context.CancelFunc) (*Pool, error) {
	res, err := cli.Get(context.Background(), sName, clientv3.WithPrefix())
	if err != nil {
		cancel()
		return nil, err
	}
	pool := &sync.Map{}
	// 将获取到的ip和port保存到本地的map
	for _, kv := range res.Kvs {
		pool.Store(string(kv.Key), string(kv.Value))
	}

	return &Pool{
		Ctx:         ctx,
		CtxCancel:   cancel,
		EtcdClient:  cli,
		ServiceName: sName,
		IpPool:      pool,
	}, nil
}

func (e *Pool) watcher() {
	watchChan := e.EtcdClient.Watch(context.Background(), e.ServiceName, clientv3.WithPrefix())
	for {
		select {
		case val := <-watchChan:
			for _, event := range val.Events {
				switch event.Type {
				case 0: // 0 是有数据增加
					e.IpPool.Store(string(event.Kv.Key), string(event.Kv.Value))
					//_global.GO_LOG.Info(e.scheme + "服务数量增加")
				case 1: // 1是有数据减少
					e.IpPool.Delete(string(event.Kv.Key))
					//_global.GO_LOG.Info(e.scheme + "服务数量减少")
				}
			}
		case <-e.Ctx.Done():
			return
		}
	}
}

func (e *Pool) Height() int {
	h := 0
	e.IpPool.Range(func(key, value interface{}) bool {
		h++
		return true
	})
	return h
}

func (e *Pool) WeightedRoundRobin() string {
	ipPool := make([]PoolData, 0)
	e.IpPool.Range(func(key, value interface{}) bool {
		tA, ok := value.(string)
		if !ok {
			return false
		}
		data := rpcConf.Config{}
		err := json.Unmarshal([]byte(tA), &data)
		if err != nil {
			return false
		}
		ipPool = append(ipPool, PoolData{
			Address: data.Ip + ":" + strconv.Itoa(int(data.Port)),
			Weight:  int(data.Weight),
			Current: 0,
		})
		return true
	})
	if len(ipPool) == 0 {
		return ""
	}
	round := NewWeightedRoundRobin(ipPool)
	return round.Next()
}
