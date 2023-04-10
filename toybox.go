package toybox

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/inoth/toybox/components"
	"github.com/inoth/toybox/services"
)

type Register struct {
	cmps []components.Components
	svcs []services.Service
	once sync.Once
}

func New(cmps ...components.Components) *Register {
	if len(cmps) <= 0 {
		fmt.Println("Warn: No components have been loaded yet.")
	}
	reg := &Register{
		cmps: cmps,
	}
	reg.once.Do(func() {
		for _, cmp := range reg.cmps {
			must(cmp.Init())
		}
	})
	return reg
}

func (reg *Register) Run(svcs ...services.Service) {
	reg.svcs = svcs
	for _, svc := range reg.svcs {
		go func(ctx context.Context, service services.Service) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}()
			must(service.Start())
		}(context.Background(), svc)
	}

	// 监听退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	fmt.Printf("Received signal %s, exiting...\n", sig)
}

func must(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
