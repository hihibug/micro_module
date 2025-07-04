package nacos

import (
	"fmt"
	rpc "github.com/hihibug/microdule/v2/Framework/rpc/config"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
	"strconv"
	"strings"
)

type Nacos struct {
	Srv               naming_client.INamingClient
	SrvConf           []constant.ServerConfig
	Conf              Config
	SubscribeCallback func([]rpc.Config)
}

func NewNacos(conf Config) *Nacos {
	sc := strings.Split(conf.Addr, ";")
	serverConfigs := make([]constant.ServerConfig, len(sc))
	for k, v := range sc {
		c := strings.Split(v, ":")
		serverConfigs[k] = constant.ServerConfig{
			IpAddr: c[0],
		}
		if len(c) > 1 {
			port, _ := strconv.Atoi(c[1])
			serverConfigs[k].Port = uint64(port)

			ports := strings.Split(c[1], "-")
			if len(ports) > 1 {
				grpcPort, _ := strconv.Atoi(c[2])
				serverConfigs[k].GrpcPort = uint64(grpcPort)
			}
		}
	}

	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  conf.ClientConfig,
	})
	if err != nil {
		panic("registry nacos create client err :" + err.Error())
	}

	return &Nacos{
		Srv:     namingClient,
		Conf:    conf,
		SrvConf: serverConfigs,
	}
}

func (n *Nacos) Client() any {
	return n
}

func (n *Nacos) RegisterService(conf *rpc.Config) error {
	n.Conf.RegisterData = conf
	regConf := vo.RegisterInstanceParam{
		Ip:          conf.Ip,
		Port:        conf.Port,
		ServiceName: conf.ServiceName,
		Weight:      float64(conf.Weight),
		Enable:      true,
		Healthy:     true,
		Ephemeral:   n.Conf.Ephemeral,
		Metadata:    conf.Metadata,
		ClusterName: conf.ClusterName, // 默认值DEFAULT
		GroupName:   conf.GroupName,   // 默认值DEFAULT_GROUP
	}

	_, err := n.Srv.RegisterInstance(regConf)
	if n.SubscribeCallback != nil {
		err := n.Srv.Subscribe(&vo.SubscribeParam{
			ServiceName: n.Conf.RegisterData.ServiceName,
			GroupName:   n.Conf.RegisterData.GroupName,             // 默认值DEFAULT_GROUP
			Clusters:    []string{n.Conf.RegisterData.ClusterName}, // 默认值DEFAULT
			SubscribeCallback: func(services []model.Instance, err error) {
				if err != nil {
					log.Printf("nacos subscribe err: %s", err.Error())
				}
				result := make([]rpc.Config, len(services))
				for _, service := range services {
					result = append(result, rpc.Config{
						Ip:          service.Ip,
						Port:        service.Port,
						ServiceName: service.ServiceName,
						Weight:      uint64(service.Weight),
						GroupName:   n.Conf.RegisterData.GroupName,
						ClusterName: n.Conf.RegisterData.ClusterName,
					})
				}
				n.SubscribeCallback(result)
			},
		})
		if err != nil {
			return err
		}
	}

	log.Printf("rpc register: %s", "nacos")
	return err
}

func (n *Nacos) DeregisterService() error {
	conf := vo.DeregisterInstanceParam{
		Ip:          n.Conf.RegisterData.Ip,
		Port:        n.Conf.RegisterData.Port,
		ServiceName: n.Conf.RegisterData.ServiceName,
		Ephemeral:   n.Conf.Ephemeral,
		Cluster:     n.Conf.RegisterData.ClusterName, // 默认值DEFAULT
		GroupName:   n.Conf.RegisterData.GroupName,   // 默认值DEFAULT_GROUP
	}

	_, err := n.Srv.DeregisterInstance(conf)
	if n.SubscribeCallback != nil {
		n.UnsubscribeService()
	}
	n.Srv.CloseClient()
	return err
}

func (n *Nacos) DiscoveryService(sName, GroupName, ClusterName string) (string, error) {
	cluster := make([]string, 0)
	if ClusterName != "" {
		cluster = append(cluster, ClusterName)
	}
	instances, err := n.Srv.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: sName,
		GroupName:   GroupName,
		Clusters:    cluster,
	})
	if err != nil {
		return "", err
	}

	serverAddr := fmt.Sprintf("%s:%d", instances.Ip, instances.Port)
	return serverAddr, nil
}

func (n *Nacos) SubscribeService(callback func([]rpc.Config)) {
	n.SubscribeCallback = callback
}

func (n *Nacos) UnsubscribeService() {
	err := n.Srv.Unsubscribe(&vo.SubscribeParam{
		ServiceName: n.Conf.RegisterData.ServiceName,
		GroupName:   n.Conf.RegisterData.GroupName,             // 默认值DEFAULT_GROUP
		Clusters:    []string{n.Conf.RegisterData.ClusterName}, // 默认值DEFAULT
	})
	if err != nil {
		log.Printf("nacos unsubscribe err: %s", err.Error())
	}
}
