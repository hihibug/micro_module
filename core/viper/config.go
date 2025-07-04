package viper

import (
	framework "github.com/hihibug/microdule/v2/Framework/config"
	db "github.com/hihibug/microdule/v2/Framework/db/config"
	rest "github.com/hihibug/microdule/v2/Framework/http/config"
	rpc "github.com/hihibug/microdule/v2/Framework/rpc/config"
	"github.com/hihibug/microdule/v2/core/zap"
)

type Config struct {
	DB        *db.Config        `json:"db" yaml:"db"`
	Log       *zap.Config       `json:"log" yaml:"log"`
	Rest      *rest.Config      `json:"http" yaml:"http"`
	Rpc       *rpc.Config       `json:"rpc" yaml:"rpc"`
	Framework *framework.Config `json:"Framework" yaml:"Framework"`
}
