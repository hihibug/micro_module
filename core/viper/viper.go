package viper

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/hihibug/micro_module/Framework/config"
	"github.com/hihibug/micro_module/core/gorm"
	"github.com/spf13/viper"
)

type ConfigViper interface {
	ConfigToGormMysql(opt gorm.OptConfig) *gorm.Config
}

type Viper struct {
	Vp   *viper.Viper
	Data *config.ConfigData
}

func NewViper(path string, defaultConf func(*viper.Viper)) *Viper {
	conf := &config.ConfigData{}
	v := viper.New()
	defaultConf(v)
	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		log.Println("config read err:", err)
		return &Viper{nil, conf}
	}

	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		_ = v.Unmarshal(conf)
	})

	err = v.Unmarshal(conf)
	if err != nil {
		log.Println("config unmarshal err:", err)
		return &Viper{nil, conf}
	}

	return &Viper{v, conf}
}

func (v *Viper) GetVal(name string) any {
	return v.Vp.Get(name)
}

func (v *Viper) GetBindVal() *config.ConfigData {
	return v.Data
}

func (v *Viper) Client() *viper.Viper {
	return v.Vp
}

// func (v *Viper) ConfigToGormMysql(opt gorm.OptConfig) *gorm.Config {
// 	return &gorm.Config{
// 		DbType:      v.Data.DB.DbType,
// 		MaxIdleCons: v.Data.DB.MaxIdleCons,
// 		MaxOpenCons: v.Data.DB.MaxOpenCons,
// 		LogMode:     v.Data.DB.LogMode,
// 		Opt:         opt,
// 		Mysql: &gorm.MysqlConfig{
// 			Path:     v.Data.DB.Path,
// 			Config:   v.Data.DB.Config,
// 			Dbname:   v.Data.DB.Dbname,
// 			Username: v.Data.DB.Username,
// 			Password: v.Data.DB.Password,
// 		},
// 	}
// }
