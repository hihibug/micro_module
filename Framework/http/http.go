package http

import (
	"context"

	http "github.com/hihibug/micro_module/Framework/http/config"
	"github.com/hihibug/micro_module/Framework/http/request"
	"github.com/hihibug/micro_module/Framework/http/response"

	"github.com/hihibug/micro_module/Framework"
	"github.com/jedib0t/go-pretty/v6/table"
)

type (
	Http interface {
		Name() string
		Config() http.Config
		Client() any
		Run() error
		Request(c any) request.Request
		Response(c any) response.Response
	}

	MicroHttp struct {
		name    string
		Context context.Context
		client  Http
	}
)

func NewHttp(http Http) Http {
	return http
}

func NewMicroHttp(http func(*http.Config) Http) Framework.Option {
	return func(options *Framework.Options) Framework.OptionHandle {
		return &MicroHttp{
			name:    "http",
			Context: context.Background(),
			client:  http(options.Config.Data.Http),
		}
	}
}

func (h *MicroHttp) Name() string {
	return h.name
}

func (h *MicroHttp) Start(runInfo chan table.Row) error {
	runInfo <- table.Row{h.name, h.client.Name(), "http://127.0.0.1:" + h.client.Config().Port}
	return h.client.Run()
}

func (h *MicroHttp) Close() error {
	return nil
}

func (h *MicroHttp) Client() any {
	return h.client
}
