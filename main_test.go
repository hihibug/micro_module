package main_test

import (
	"fmt"
	"github.com/hihibug/microdule/v2/Framework/http"
	httpConf "github.com/hihibug/microdule/v2/Framework/http/config"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hihibug/microdule/v2/Framework"
	microdulefiber "github.com/hihibug/microdule/v2/core/fiber"
)

type global struct {
	Http http.Http
}

var g *global

func TestNewService(t *testing.T) {
	//初始化服务 config初始化键值为0
	s := Framework.NewService(
		Framework.Config("core/config.yml"),
		Framework.Name("test"),
	)

	//初始化组件
	s.Init(
		http.NewMicroHttp(microdulefiber.NewFiber(&httpConf.Config{
			LogPath: "",
			UseHtml: false,
			Addr:    "8999",
		})),
	)

	g = &global{
		s.Microdule("http").Client().(*microdulefiber.Fiber),
	}

	r := s.Microdule("http").Client().(*microdulefiber.Fiber).Route

	r.All("/", func(ctx *fiber.Ctx) error {
		fmt.Println(1)
		return g.Http.Response(ctx).Ok()
	})

	fmt.Println(s.Start())
}
