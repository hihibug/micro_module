package http

import (
	"context"
	"github.com/hihibug/microdule/v2/Framework/http/config"
	"github.com/hihibug/microdule/v2/Framework/http/request"
	"github.com/hihibug/microdule/v2/Framework/http/response"

	"github.com/hihibug/microdule/v2/Framework"
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

func NewMicroHttp(http Http) Framework.Option {
	return func(options *Framework.Options) Framework.OptionHandle {
		return &MicroHttp{
			name:    "http",
			Context: context.Background(),
			client:  http,
		}
	}
}

func (h *MicroHttp) Name() string {
	return h.name
}

func (h *MicroHttp) Start(runInfo chan table.Row) error {
	runInfo <- table.Row{h.name, h.client.Name(), "http://127.0.0.1:" + h.client.Config().Addr}
	return h.client.Run()

	// randomNum := rand.Intn(10000)
	// fmt.Println(randomNum)
	// time.Sleep(time.Duration(randomNum) * time.Millisecond)
	// panic("123123")
}

func (h *MicroHttp) Close() error {
	return nil
}

func (h *MicroHttp) Client() any {
	return h.client
}
