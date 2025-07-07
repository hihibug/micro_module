package micro_module_gin

import (
	"io"
	"net/http"
	"os"
	"time"

	httpConf "github.com/hihibug/micro_module/Framework/http/config"
	"github.com/hihibug/micro_module/Framework/http/request"
	"github.com/hihibug/micro_module/Framework/http/response"
	"github.com/hihibug/micro_module/Framework/http/validator"

	"github.com/gin-gonic/gin"
	"github.com/hihibug/micro_module/core/gin/middleware"
	"github.com/hihibug/micro_module/core/utils"
)

type Gin struct {
	name      string
	Route     *gin.Engine
	conf      *httpConf.Config
	Validator validator.Validator
}

func NewGin(conf *httpConf.Config) *Gin {
	gin.SetMode(gin.ReleaseMode)

	// 日志写入
	defPath, _ := os.Getwd()
	path := defPath + "/" + conf.LogPath
	if ok, _ := utils.PathExists(path); !ok { // 判断是否有Director文件夹
		_ = os.Mkdir(path, os.ModePerm)
	}
	accessLogPath := path + "/access-" + time.Now().Format("2006-01-02") + ".log"
	// 记录到文件。
	f, _ := os.Create(accessLogPath)
	gin.DefaultWriter = io.MultiWriter(f)

	var route = gin.Default()

	// 初始化页面
	if conf.UseHtml {
		defPath, _ := os.Getwd()
		route.Delims(conf.DelimsLeft, conf.DelimsRight)
		route.Static(defPath+"/"+conf.StaticPath, defPath+"/"+conf.HtmlPath)
		route.LoadHTMLGlob(defPath + "/" + conf.HtmlPath + "/*")
	}

	//注册GinCors
	route.Use(middleware.GinCors(), middleware.GinErrorHttp)
	route.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "page not found", "code": http.StatusNotFound})
	})

	// 验证器

	return &Gin{
		"Gin",
		route,
		conf,
		validator.NewValidator("zh"),
	}
}

func (g *Gin) Name() string {
	return g.name
}

func (g *Gin) Config() httpConf.Config {
	return *g.conf
}

func (g *Gin) Client() any {
	return g
}

func (g *Gin) Run() error {
	s := &http.Server{
		Addr:           ":" + g.conf.Port,
		Handler:        g.Route,
		ReadTimeout:    time.Duration(g.conf.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(g.conf.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s.ListenAndServe()
}

func (g *Gin) Response(c any) response.Response {
	return NewGinResponse(c.(*gin.Context))
}

func (g *Gin) Request(c any) request.Request {
	return NewGinRequest(c.(*gin.Context), g.Validator)
}
