package Framework_test

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hihibug/micro_module/Framework/http"

	"github.com/hihibug/micro_module/Framework"
	micro_modulefiber "github.com/hihibug/micro_module/core/fiber"
)

type global struct {
	Http http.Http
}

var g *global

func initGlobal(s Framework.Service) {
	g = &global{
		Http: s.Modules("http").Client().(http.Http),
	}
}

func TestNewService(t *testing.T) {
	//初始化服务 config初始化键值为0
	s := Framework.NewService("")

	//初始化组件
	s.Init(
		http.NewMicroHttp(micro_modulefiber.NewFiber),
	)

	initGlobal(s)

	fmt.Println(g.Http.Name())

	g.Http.Client().(*micro_modulefiber.Fiber).Route.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World 👋!")
	})

	fmt.Println(s.Start())
}

// func GoMysql(num, cnum int) {
// 	var wg sync.WaitGroup
// 	ch := make(chan struct{}, cnum)
// 	for i := 0; i < num; i++ {
// 		ch <- struct{}{}
// 		wg.Add(1)
// 		go func(i int) {
// 			defer wg.Done()
// 			a := make([]map[string]interface{}, 0)
// 			err := global.Gorm.Client().Table("users").Find(&a).Error
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			log.Println(a)
// 			time.Sleep(time.Second)
// 			<-ch
// 		}(i)
// 	}
// 	wg.Wait()
// }

// func GoRedis(num, cnum int) {
// 	var wg sync.WaitGroup
// 	ch := make(chan struct{}, cnum)
// 	for i := 0; i < num; i++ {
// 		ch <- struct{}{}
// 		wg.Add(1)
// 		go func(i int) {
// 			defer wg.Done()
// 			c := global.Redis.Client().Get("test-11")
// 			if c.Err() != nil {
// 				log.Println(c.Err())
// 			}
// 			log.Println(c.Val())
// 			time.Sleep(time.Second)
// 			<-ch
// 		}(i)
// 	}
// 	wg.Wait()
// }
