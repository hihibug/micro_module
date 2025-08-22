package Framework

import (
	"context"

	"github.com/hihibug/micro_module/Framework/config"
	"github.com/hihibug/micro_module/Framework/log"
	"github.com/hihibug/micro_module/core/viper"
	"github.com/jedib0t/go-pretty/v6/table"
)

type (
	OptionHandle interface {
		Name() string
		Start(chan table.Row) error
		Close() error
		Client() any
	}
	Option func(*Options) OptionHandle

	Options struct {
		Name    string
		Context context.Context
		Config  config.Config
		Log     log.Log
	}

	service struct {
		opts   Options
		extend map[string]OptionHandle
	}
)

// Run implements Service.
func (*service) Run() error {
	panic("unimplemented")
}

func newService(logPath string) *service {
	opt := Options{
		Context: context.Background(),
	}
	if logPath == "" {
		logPath = "./config.yaml"
	}
	newConfig(logPath)(&opt)
	opt.Name = opt.Config.GetBindVal().Name
	opt.Log = log.NewLog(opt.Config.GetBindVal().Log)

	return &service{
		opts:   opt,
		extend: make(map[string]OptionHandle),
	}
}

func newConfig(path string) Option {
	return func(options *Options) OptionHandle {
		options.Config = viper.NewViper(path, config.InitConfig)
		return nil
	}
}
