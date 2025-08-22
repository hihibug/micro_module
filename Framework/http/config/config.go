package http

import (
	"github.com/spf13/viper"
)

var prfix = "Http."

type Config struct {
	IsLog          bool   `json:"IsLog" yaml:"IsLog"`                   // 访问日志是否开启
	LogPath        string `json:"LogPath" yaml:"LogPath"`               // 访问日志地址
	UseHtml        bool   `json:"UseHtml" yaml:"UseHtml"`               // html 页面开关
	StaticPath     string `json:"StaticPath" yaml:"StaticPath"`         // 静态文件地址
	HtmlPath       string `json:"HtmlPath" yaml:"HtmlPath"`             // html模板地址
	DelimsLeft     string `json:"Delimsleft" yaml:"Delimsleft"`         // html模板标签左
	DelimsRight    string `json:"DelimsRight" yaml:"DelimsRight"`       // html模板标签右
	Port           string `json:"Port" yaml:"Port"`                     // 端口
	BindIP         string `json:"BindIp" yaml:"BindIp"`                 // 启动地址
	ReadTimeout    int    `json:"ReadTimeout" yaml:"ReadTimeout"`       // 读超时
	WriteTimeout   int    `json:"WriteTimeout" yaml:"WriteTimeout"`     // 写超时
	MaxHeaderBytes int    `json:"MaxHeaderBytes" yaml:"MaxHeaderBytes"` // 最大请求头
}

func InitConfig(v *viper.Viper) {
	v.SetDefault(prfix+"Port", "8080")
	v.SetDefault(prfix+"BindIp", "127.0.0.1")
	v.SetDefault(prfix+"ReadTimeout", 10)
	v.SetDefault(prfix+"WriteTimeout", 10)
	// v.Vp.SetDefault(prfix+"MaxHeaderBytes", 1<<20)
	v.SetDefault(prfix+"IsLog", false)
	v.SetDefault(prfix+"LogPath", "storage/log")
	v.SetDefault(prfix+"UseHtml", false)
}
