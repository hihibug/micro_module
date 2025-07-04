package Framework

import (
	"errors"
	"fmt"
	"github.com/hihibug/microdule/v2/Framework/_global"
	"log"
	"os"
	"runtime/debug"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/exp/maps"
)

type Service interface {
	Name() string
	Init(...Option)
	Options() *Options
	Microdule(string) OptionHandle
	Close()
	Start() error
}

func NewService(opt ...Option) Service {
	_global.CliTable = table.NewWriter()
	_global.CliTable.SetOutputMirror(os.Stdout)

	return newService(opt...)
}

func (s *service) Name() string {
	return s.opts.Name
}

func (s *service) Init(opts ...Option) {
	if s.opts.Name == "" {
		panic("name is empty")
	}
	_global.CliTable.SetTitle(s.opts.Name + " framework start")
	_global.CliTable.AppendHeader(table.Row{"Modules", "Components", ""})
	for _, o := range opts {
		t := o(&s.opts)
		maps.Copy(s.extend, map[string]OptionHandle{t.Name(): t})
	}
}

func (s *service) Options() *Options {
	return &s.opts
}

func (s *service) Microdule(name string) OptionHandle {
	return s.extend[name]
}

func (s *service) Close() {
	for _, t := range s.extend {
		t.Close()
	}
}

func (s *service) Start() error {
	sucNum := 0
	infoChan := make(chan table.Row, len(s.extend))
	stopChan := make(chan string)
	for k, f := range s.extend {
		go func(handler func(chan table.Row) error, k string) {
			defer func() {
				// 每个协程内部使用recover捕获可能在调用逻辑中发生的panic
				e := recover()
				if e != nil {
					log.Println(e)
					//打印错误堆栈信息
					debug.PrintStack()
					stopChan <- fmt.Sprintf("%s", e)
				}
			}()

			err := handler(infoChan)
			if err != nil {
				stopChan <- err.Error()
			}
		}(f.Start, k)
	}

	for {
		select {
		case err := <-stopChan:
			return errors.New(err)
		case info := <-infoChan:
			sucNum += 1
			_global.CliTable.AppendRow(info)
			_global.CliTable.AppendSeparator()
			if sucNum == len(s.extend) {
				_global.CliTable.SetCaption("All modules are running.\n")
				_global.CliTable.SetStyle(table.StyleLight)
				_global.CliTable.Render()
			}
		}
	}
}
