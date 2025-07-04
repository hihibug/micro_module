package Framework

import (
	"context"

	"reflect"

	"github.com/hihibug/microdule/v2/core/viper"
	"github.com/hihibug/microdule/v2/core/zap"
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
)

func newOptions(opts ...Option) Options {
	opt := Options{
		Context: context.Background(),
	}

	for k, o := range opts {
		if k > 0 && reflect.DeepEqual(opt.Config, viper.Viper{}) {
			opt.Config = viper.NewViper("config.yml")
		}
		if opt.Config.Err != nil {
			panic(opt.Config.Err)
		}
		o(&opt)
	}
	if opt.Log == nil {
		opt.Log = zap.NewZap(opt.Config.Data.Log)
	}

	if opt.Name == "" {
		opt.Name = opt.Config.Data.Framework.Name
	}

	return opt
}

func Name(n string) Option {
	return func(options *Options) OptionHandle {
		options.Name = n
		return nil
	}
}

func Config(path string) Option {
	return func(options *Options) OptionHandle {
		options.Config = viper.NewViper(path)
		return nil
	}
}

//	func Gorm(dbConf *gorm.Config) Option {
//		return func(options *Options) {
//			if dbConf == nil {
//				dbConf = options.Config.ConfigToGormMysql(gorm.SetGormConfig(gorm.GetGormConfigStruct()))
//			}
//			db, err := gorm.NewGorm(dbConf)
//			if err != nil {
//				panic("mysql error " + err.Error())
//			}
//			options.Gorm = db
//		}
//	}
//
//	func Etcd(e *etcd.Config) Option {
//		return func(options *Options) {
//			if e == nil {
//				e = options.Config.Data.Etcd
//			}
//			etd, err := etcd.NewEtcd(e)
//			if err != nil {
//				panic("etcd error " + err.Error())
//			}
//			options.Etcd = etd
//		}
//	}
//
//	func Redis(r *redis.Config) Option {
//		return func(options *Options) {
//			if r == nil {
//				r = options.Config.Data.Redis
//			}
//			rds, err := redis.NewRedis(r)
//			if err != nil {
//				panic("redis error " + err.Error())
//			}
//			options.Redis = rds
//		}
//	}

//
//func Rpc(r rpc.Rpc) Option {
//	return func(options *Options) {
//		options.Rpc = r
//	}
//}
