package toybox

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/inoth/toybox/config"
	"github.com/inoth/toybox/util"
	"golang.org/x/sync/errgroup"
)

type ToyBox struct {
	option

	ctx    context.Context
	cancel context.CancelFunc
}

func New(opts ...Option) *ToyBox {
	o := option{
		id:      util.UUID(),
		version: util.UUID(),
		ctx:     context.Background(),
		sigs:    []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
	for _, opt := range opts {
		opt(&o)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &ToyBox{
		option: o,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (ce *ToyBox) ID() string      { return ce.id }
func (ce *ToyBox) Name() string    { return ce.name }
func (ce *ToyBox) Version() string { return ce.version }

func (ce *ToyBox) Run() (err error) {
	fmt.Printf("server start %s\n", ce.ID())

	if ce.cfg == nil {
		return ErrNotConfig
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, ce.sigs...)

	wg := sync.WaitGroup{}
	eg, ctx := errgroup.WithContext(ce.ctx)

	for _, svc := range ce.svcs {
		svc := svc
		cm, ok := svc.(config.ConfigureMatcher)
		if !ok {
			continue
		}
		if err := ce.cfg.PrimitiveDecode(svc.(config.ConfigureMatcher)); err != nil {
			return err
		}
		eg.Go(func() error {
			<-ctx.Done()
			fmt.Printf("Done %s ...............\n", cm.Name())
			return svc.Stop(ctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			fmt.Printf("Start %s ...............\n", cm.Name())
			return svc.Start(ctx)
		})
	}

	wg.Wait()

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			fmt.Printf("Done server %s ...............\n", ce.ID())
			return nil
		case <-c:
			fmt.Printf("Done server %s ...............\n", ce.ID())
			return ce.Stop()
		}
	})
	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (ce *ToyBox) Stop() error {
	if ce.cancel != nil {
		ce.cancel()
	}
	return nil
}
