package http_test

//
import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hihibug/micro_module/Framework/http"
	httpConf "github.com/hihibug/micro_module/Framework/http/config"
	micro_module_fiber "github.com/hihibug/micro_module/core/fiber"
	micro_module_gin "github.com/hihibug/micro_module/core/gin"

	//"github.com/hihibug/Framework/v2/Http/rest_gin"
	"testing"
)

func TestFiberRest(t *testing.T) {
	r := http.NewHttp(micro_module_fiber.NewFiber(&httpConf.Config{
		LogPath: "",
		UseHtml: false,
		Port:    ":8080",
	}))

	if r.Client() != nil {
		rs := r.Client().(*micro_module_fiber.Fiber)
		fmt.Println(rs)
		err := rs.Validator.RegisterValidator("test", "不能为空且不等于admin", notNullAndAdmin)
		fmt.Println(err)
		a := rs.Route.Group("")
		{
			a.All("/test", func(c *fiber.Ctx) error {
				type aa struct {
					Name string `form:"name" json:"name" validate:"test,min=5,max=20"`
					Age  uint   `form:"age" json:"age" validate:"lte=60,gte=0"`
				}
				var a aa
				err := r.Request(c).GetVerifyVal(&a)
				if err != nil {
					fmt.Println(err)
					return r.Response(c).FailWithMessage(err.Error())
				}
				fmt.Println(err)
				fmt.Println(a)
				return r.Response(c).Ok()
			})
		}

	}

	fmt.Println(r.Run())
}

func notNullAndAdmin(c validator.FieldLevel) bool {
	value := c.Field().String()
	//字段不能为空，并且不等于admin
	return value != "" && !(value == "admin")
}

func TestGinRest(t *testing.T) {
	//r := Http.NewRest(&httpConf.Config{
	//	Mode:    "gin",
	//	LogPath: "",
	//	UseHtml: false,
	//	Addr:    "8999",
	//})
	r := micro_module_gin.NewGin(&httpConf.Config{
		LogPath: "",
		UseHtml: false,
		Port:    "8999",
	})

	if r.Client() != nil {
		rs := r.Client().(*micro_module_gin.Gin)
		fmt.Println(rs)
		err := rs.Validator.RegisterValidator("test", "不能为空且不等于admin", notNullAndAdmin)
		fmt.Println(err)
		a := rs.Route.Group("")
		{
			a.Any("/test", func(c *gin.Context) {
				type aa struct {
					Name string `form:"name" json:"name" validate:"test,min=5,max=20"`
					Age  uint   `form:"age" json:"age" validate:"lte=60,gte=0"`
				}
				var a aa
				err := r.Request(c).GetVerifyVal(&a)
				fmt.Println(err)
				//.Validator.FetchValidatorError(c)
				fmt.Println(a)

			})
		}

	}

	fmt.Println(r.Run())
}
