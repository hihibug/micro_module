package nacos

import (
	"github.com/hihibug/microdule/v2/Framework/rpc/config"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
)

type Config struct {
	Addr         string `yaml:"addr" json:"addr"` // 127.0.0.1:90-1090;xxx.x.x.x:xx-xx ip:port-grpc_port;
	Ephemeral    bool   `yaml:"ephemeral" json:"ephemeral"`
	ClientConfig constant.ClientConfig
	RegisterData *config.Config
}
