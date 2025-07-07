package http

import (
	"github.com/google/wire"
	"github.com/hihibug/microdule"
	"github.com/hihibug/microdule/web"
)

var ProviderHttp = wire.NewSet(InitHttp)

// InitRedis 初始化redis
func InitHttp(global *Global) *web.Gin {
	httpData := web.NewGin(global.Config.Http).Client().(*web.Gin)
	config.RegisterRoute(httpData)
	global.Srv.Init(microdule.Http(httpData))
	return httpData
}
