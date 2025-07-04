package Framework

import (
	"fmt"
	"github.com/hihibug/microdule/v2/Framework/http"
	httpConf "github.com/hihibug/microdule/v2/Framework/http/config"
	//"github.com/hihibug/microdule_gin"
	"log"
	"sync"
	"testing"
	"time"
)

var global *Options

func TestNewService(t *testing.T) {
	//初始化服务 config初始化键值为0
	s := NewService(
		Config("core/config.yml"),
		Name("test"),
	)

	global = s.Options()

	// fmt.Println(_global.Config.Vp.Get("QiNu.Key"))
	////NewZapWriter 对log.New函数的再次封装，从而实现是否通过zap打印日志
	//gormConf := gorm.GetGormConfigStruct()
	//gorm.LogGorm(
	//	_global.Config.Data.DB.LogMode,
	//	gormConf,
	//	gorm.SetGormLogZap(zap.NewZapWriter(_global.Log.Client())),
	//)
	//
	//// 设置etcd 日志
	//_global.Config.Data.Etcd.Log = _global.Log.Client()
	//
	////初始化组件
	s.Init(
		MicroInit(http.NewMicroHttp(microdule_gin.NewGin(&httpConf.Config{
			Mode:    "fiber",
			LogPath: "",
			UseHtml: false,
			Addr:    "8999",
		}))),
		//Redis(nil),
		//Gorm(_global.Config.ConfigToGormMysql(gorm.SetGormConfig(gormConf))),
		//Etcd(_global.Config.Data.Etcd),
		//Http(microdule_gin.NewGin(&httpConf.Config{
		//	Mode:    "fiber",
		//	LogPath: "",
		//	UseHtml: false,
		//	Addr:    "8999",
		//})),

	)
	fmt.Println(global.Config.Data.Http)
	//fmt.Println()
	fmt.Println(s.Options().Http)
	//
	////关闭链接
	//defer s.Close()
	//
	////开启rest
	//http := s.Http().Client()
	//
	//a := http.Route.Group("")
	//{
	//	a.GET("/test", func(context *rest_gin.Context) {
	//		fmt.Println("test")
	//		_global.Log.Client().Info("test")
	//	})
	//	a.GET("/err", func(c *rest_gin.Context) {
	//		panic("test")
	//	})
	//}

	//ip, _ := utils.ExternalIP()
	//_global.Config.Data.Rpc.IP = ip
	//grpcData := rpc.NewGrpc(_global.Config.Data.Rpc)
	//// config.RegisterGrpc(grpcData.Client().RpcSrv)
	//
	//register, err := grpcData.Register(_global.Etcd.Clients())
	//if err != nil {
	//	panic(err)
	//}
	//grpcData.Client().(*rpc.Grpc).EtcdRegister = register
	//s.Init(Rpc(grpcData))
	// rpc := s.Rpc().Client().Grpc
	//register, err := s.Rpc().Client().Register(_global.Etcd.Clients())
	//if err != nil {
	//	os.Exit(0)
	//}
	//go register.ListenLeaseRespChan()
	//
	//grpcData.Run()
	s.Options().Http.Run()
}

func GoMysql(num, cnum int) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, cnum)
	for i := 0; i < num; i++ {
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			a := make([]map[string]interface{}, 0)
			err := global.Gorm.Client().Table("users").Find(&a).Error
			if err != nil {
				log.Println(err)
			}
			log.Println(a)
			time.Sleep(time.Second)
			<-ch
		}(i)
	}
	wg.Wait()
}

func GoRedis(num, cnum int) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, cnum)
	for i := 0; i < num; i++ {
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c := global.Redis.Client().Get("test-11")
			if c.Err() != nil {
				log.Println(c.Err())
			}
			log.Println(c.Val())
			time.Sleep(time.Second)
			<-ch
		}(i)
	}
	wg.Wait()
}
