package etcd

import (
	"github.com/hihibug/micro_module/Framework/rpc/config"
	"github.com/hihibug/micro_module/core/etcd"
)

type Config struct {
	Etcd         etcd.Config
	RegisterData *config.Config
}
