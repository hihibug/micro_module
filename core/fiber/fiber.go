package micro_module_fiber

import (
	"os"
	"time"

	"github.com/hihibug/micro_module/Framework/http"
	httpConf "github.com/hihibug/micro_module/Framework/http/config"
	"github.com/hihibug/micro_module/Framework/http/request"
	"github.com/hihibug/micro_module/Framework/http/response"
	"github.com/hihibug/micro_module/Framework/http/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/hihibug/micro_module/core/utils"
)

type Fiber struct {
	name      string
	Route     *fiber.App
	conf      *httpConf.Config
	Validator validator.Validator
}

func NewFiber(conf *httpConf.Config) http.Http {
	defPath, _ := os.Getwd()
	fc := fiber.Config{
		DisableStartupMessage: true,
	}
	if conf.UseHtml {
		fc.Views = html.New(defPath+"/"+conf.HtmlPath, ".html")
	}

	app := fiber.New(fc)

	// 请求日志
	path := defPath + "/" + conf.LogPath
	if ok, _ := utils.PathExists(path); !ok { // 判断是否有Director文件夹
		_ = os.Mkdir(path, os.ModePerm)
	}
	accessLogPath := path + "/access-" + time.Now().Format("2006-01-02") + ".log"
	// 记录到文件。
	f, _ := os.Create(accessLogPath)
	app.Use(logger.New(logger.Config{
		Output:     f,
		Format:     "[Fiber] ${time} | ${status} |  ${latency} | ${ip} | ${method}  ${path} \n",
		TimeFormat: "2006/01/02 15:04:05",
		TimeZone:   "Asia/Shanghai",
	}))

	// 初始化页面
	if conf.UseHtml {
		app.Static(defPath+"/"+conf.StaticPath, defPath+"/"+conf.HtmlPath)
	}

	return &Fiber{
		name:      "Fiber",
		Route:     app,
		conf:      conf,
		Validator: validator.NewValidator("zh"),
	}
}

func (f *Fiber) Name() string {
	return f.name
}

func (f *Fiber) Config() httpConf.Config {
	return *f.conf
}

func (f *Fiber) Client() any {
	return f
}

func (f *Fiber) Run() error {
	// log.Printf("http  port: %s \n", f.conf.Addr)
	return f.Route.Listen(":" + f.conf.Port)
}

func (f *Fiber) Response(c any) response.Response {
	return NewFiberResponse(c.(*fiber.Ctx))
}

func (f *Fiber) Request(c any) request.Request {
	return NewFiberRequest(c.(*fiber.Ctx), f.Validator)
}
