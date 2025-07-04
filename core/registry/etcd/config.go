package etcd

import (
	"github.com/hihibug/microdule/v2/Framework/rpc/config"
	"github.com/hihibug/microdule/v2/core/etcd"
)

type Config struct {
	Etcd         etcd.Config
	RegisterData *config.Config
}
