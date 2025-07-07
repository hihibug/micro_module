package pkg

import (
	"project_cli/internal/config"

	"github.com/google/wire"
	"github.com/hihibug/micro_module/Framework"
	"github.com/hihibug/micro_module/Framework/http"
	micro_module_fiber "github.com/hihibug/micro_module/core/fiber"
)

var ProviderPkg = wire.NewSet(InitPkg)

type Pkg struct {
	Srv  Framework.Service
	Http http.Http
}

func InitPkg() *Pkg {
	pkg := &Pkg{}
	pkg.Srv = Framework.NewService("./config.yaml")

	pkg.Srv.Init(
		http.NewMicroHttp(micro_module_fiber.NewFiber),
	)

	pkg.Http = pkg.Srv.Modules("http").(http.Http)
	config.RegisterRouter(pkg.Http.Client().(*micro_module_fiber.Fiber).Route)

	return pkg
}
