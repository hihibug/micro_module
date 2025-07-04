package task

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"sync"

	"github.com/hihibug/microdule/v2/Framework/_global"
	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/exp/maps"
)

type Coroutine interface {
	Register(name string, handle func(*sync.WaitGroup) error) CoroutineClose
	WorkNum() int
	Start() error
	Close()
}

type CoroutineClose interface {
	HandleClose(func())
}

type CoroutineCloseStruct struct {
	T *CoroutineStruct
}

type CoroutineStruct struct {
	Ctx          context.Context
	Handles      map[string]func(*sync.WaitGroup) error
	HandlesClose []func()
}

func NewCoroutine() Coroutine {
	return &CoroutineStruct{
		Ctx:          context.Background(),
		Handles:      make(map[string]func(*sync.WaitGroup) error),
		HandlesClose: make([]func(), 0),
	}
}

func (t *CoroutineStruct) Register(name string, handle func(*sync.WaitGroup) error) CoroutineClose {
	maps.Copy(t.Handles, map[string]func(*sync.WaitGroup) error{name: handle})
	return &CoroutineCloseStruct{
		T: t,
	}
}

func (t *CoroutineStruct) WorkNum() int {
	return len(t.Handles)
}

func (t *CoroutineCloseStruct) HandleClose(handle func()) {
	t.T.HandlesClose = append(t.T.HandlesClose, handle)
}

func (t *CoroutineStruct) Start() error {
	var wg sync.WaitGroup

	stopChan := make(chan string)
	// 假设我们要调用handlers这么多个服务
	for k, f := range t.Handles {
		wg.Add(1)
		// 每个函数启动一个协程
		go func(handler func(*sync.WaitGroup) error, k string) {
			defer func(wg *sync.WaitGroup) {
				wg.Done()
				// 每个协程内部使用recover捕获可能在调用逻辑中发生的panic
				e := recover()
				if e != nil {
					log.Println(e)
					//打印错误堆栈信息
					debug.PrintStack()
					stopChan <- fmt.Sprintf("%s", e)
				}
			}(&wg)
			// 取第一个报错的handler调用逻辑，并最终向外返回
			err := handler(&wg)
			if err != nil {
				wg.Done()
				stopChan <- err.Error()
			}
		}(f, k)
	}

	for {
		select {
		case err := <-stopChan:
			fmt.Println("获取错误")
			return errors.New(err)
		default:
			wg.Wait()

			_global.CliTable.SetCaption("Simple Table with 3 Rows.\n")
			_global.CliTable.AppendSeparator()
			_global.CliTable.SetStyle(table.StyleLight)
			_global.CliTable.Render()
		}
	}
	// err := <-stopChan
	// return errors.New(err)
}

func (t *CoroutineStruct) Close() {
	for _, v := range t.HandlesClose {
		v()
	}
}
