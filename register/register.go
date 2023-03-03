package register

import (
	"context"
	"fmt"

	"os"
	"sync"

	"github.com/inoth/ino-toybox/components"
	"github.com/inoth/ino-toybox/services"
)

var initOnce sync.Once

type Register struct {
	cmps []components.Components
	svcs []services.Service
}

func NewRegister(cmps ...components.Components) *Register {
	if len(cmps) <= 0 {
		fmt.Println("Warn: No components have been loaded yet.")
	}
	return &Register{
		cmps: cmps,
	}
}

func (reg *Register) Init() *Register {
	initOnce.Do(func() {
		for _, cmp := range reg.cmps {
			must(cmp.Init())
		}
	})
	return reg
}

func (reg *Register) SubStart(svcs ...services.Service) *Register {
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
	return reg
}

func (reg *Register) Start(svcs services.Service) {
	must(svcs.Start())
}

func (reg *Register) Stop() {
	for _, svcs := range reg.svcs {
		svcs.Stop()
	}
}

func must(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
