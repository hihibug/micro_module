package config

import (
	"errors"
	"reflect"

	"github.com/spf13/viper"
)

var (
	ErrEmptyDirector = errors.New("empty log director")
	ErrEmptyConfig   = errors.New("empty config")
	prfix            = "Log."
)

type Config struct {
	Level         string `json:"level" yaml:"level"`
	Format        string `json:"format" yaml:"format"`
	Prefix        string `json:"prefix" yaml:"prefix"`
	Director      string `json:"director"  yaml:"director"`
	LinkName      string `json:"linkName" yaml:"linkName"`
	ShowLine      bool   `json:"showLine" yaml:"showLine"`
	EncodeLevel   string `json:"encodeLevel" yaml:"encodeLevel"`
	StacktraceKey string `json:"stacktraceKey" yaml:"stacktraceKey"`
	LogInConsole  bool   `json:"logInConsole" yaml:"logInConsole"`
}

func (c *Config) Validate() error {
	if reflect.DeepEqual(c, &Config{}) {
		return ErrEmptyConfig
	}

	if c.Director == "" {
		return ErrEmptyDirector
	}

	return nil
}

func InitConfig(v *viper.Viper) {
	v.SetDefault(prfix+"Level", "info")
	v.SetDefault(prfix+"Format", "json")
	v.SetDefault(prfix+"Prefix", "default-log")
	v.SetDefault(prfix+"StacktraceKey", "stacktrace")
	v.SetDefault(prfix+"Director", "storage/log")
	v.SetDefault(prfix+"LogInConsole", false)
	v.SetDefault(prfix+"StacktraceKey", "stacktrace")
	v.SetDefault(prfix+"LinkName", "latest_log")
	v.SetDefault(prfix+"ShowLine", false)
}
