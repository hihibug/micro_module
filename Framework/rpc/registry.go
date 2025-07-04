package rpc

import (
	"github.com/hihibug/microdule/v2/Framework/rpc/config"
)

type Registry interface {
	Client() any
	RegisterService(conf *config.Config) error
	DiscoveryService(sName, GroupName, ClusterName string) (string, error)
	DeregisterService() error
	SubscribeService(callback func([]config.Config))
	UnsubscribeService()
}
