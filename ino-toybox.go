package inotoybox

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/inoth/ino-toybox/components"
	"github.com/inoth/ino-toybox/servers"
)

var once sync.Once

type InoToyBox struct {
	cmps []components.Component
}

func NewToyBox(cmps ...components.Component) *InoToyBox {
	if len(cmps) <= 0 {
		fmt.Println("Err: No components have been loaded yet.")
		os.Exit(1)
	}
	return &InoToyBox{
		cmps: cmps,
	}
}

// 追加组件, 满足特殊情况下判断前置组件情况下追加, 需要在 init 执行前添加
func (itb *InoToyBox) AppendComponent(cmps ...components.Component) *InoToyBox {
	if len(cmps) <= 0 {
		fmt.Println("Warn: No components have been loaded yet.")
	} else {
		itb.cmps = append(itb.cmps, cmps...)
	}
	return itb
}

func (itb *InoToyBox) Init() *InoToyBox {
	once.Do(func() {
		for _, cmp := range itb.cmps {
			must(cmp.Init())
		}
	})
	return itb
}

func (itb *InoToyBox) SubStart(svcs ...servers.Server) *InoToyBox {
	for _, svc := range svcs {
		go func(ctx context.Context, service servers.Server) {
			defer func() {
				if exception := recover(); exception != nil {
					if err, ok := exception.(error); ok {
						fmt.Printf("%v\n", err)
					} else {
						panic(exception)
					}
					os.Exit(1)
				}
			}()
			// run sub service
			must(service.Start())
		}(context.Background(), svc)
	}
	return itb
}

func (itb *InoToyBox) Start(svcs servers.Server) error {
	return svcs.Start()
}

func must(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
