package config

import (
	db "github.com/hihibug/micro_module/Framework/db/config"
	http "github.com/hihibug/micro_module/Framework/http/config"
	rpc "github.com/hihibug/micro_module/Framework/rpc/config"
	"github.com/spf13/viper"

	log "github.com/hihibug/micro_module/Framework/log/config"
)

type Config interface {
	Client() *viper.Viper
	GetBindVal() *ConfigData
	GetVal(string) any
}

type ConfigData struct {
	Name string       `json:"name" yaml:"name"`
	DB   *db.Config   `json:"db" yaml:"db"`
	Log  *log.Config  `json:"log" yaml:"log"`
	Http *http.Config `json:"http" yaml:"http"`
	Rpc  *rpc.Config  `json:"rpc" yaml:"rpc"`
}

func InitConfig(v *viper.Viper) {
	v.SetDefault("Name", "default")
	http.InitConfig(v)
	log.InitConfig(v)
}
