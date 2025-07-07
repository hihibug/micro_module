package Framework

import (
	"context"

	"github.com/hihibug/micro_module/Framework/config"
	"github.com/hihibug/micro_module/core/viper"
	"github.com/hihibug/micro_module/core/zap"
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
		Config  viper.Viper
		Log     zap.Log
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
		logPath = "config/default_conf.yaml"
	}
	Config(logPath)(&opt)

	if opt.Config.Err != nil {
		panic(opt.Config.Err)
	}

	if opt.Name == "" {
		opt.Name = opt.Config.Data.Name
	}

	if opt.Log == nil {
		opt.Log = zap.NewZap(opt.Config.Data.Log)
	}

	return &service{
		opts:   opt,
		extend: make(map[string]OptionHandle),
	}
}

func Config(path string) Option {
	return func(options *Options) OptionHandle {
		options.Config = viper.NewViper(path, config.InitConfig)
		return nil
	}
}
