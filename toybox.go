package toybox

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

func (tb *ToyBox) ID() string      { return tb.id }
func (tb *ToyBox) Name() string    { return tb.name }
func (tb *ToyBox) Version() string { return tb.version }

func (tb *ToyBox) Run() (err error) {
	log.Printf("Starting server ID:%s (PID: %d)\n", tb.ID(), os.Getpid())

	if tb.cfg == nil {
		panic(fmt.Errorf("unable to load configuration"))
	}

	watchCh := make(chan struct{}, 1)
	if watch, ok := tb.cfg.(config.Watcher); ok {
		go watch.Watche(tb.ctx, watchCh)
	}

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, tb.sigs...)

	restartCh := make(chan os.Signal, 1)
	signal.Notify(restartCh, syscall.SIGHUP)

	wg := sync.WaitGroup{}
	eg, ctx := errgroup.WithContext(tb.ctx)

	for _, svc := range tb.svcs {
		svc := svc
		cm, ok := svc.(config.ConfigureMatcher)
		if !ok {
			continue
		}
		if err := tb.cfg.PrimitiveDecode(svc.(config.ConfigureMatcher)); err != nil {
			return err
		}
		eg.Go(func() error {
			<-ctx.Done()
			log.Printf("Done %s ...............\n", cm.Name())
			return svc.Stop(ctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			log.Printf("Start %s ...............\n", cm.Name())
			return svc.Start(ctx)
		})
	}

	wg.Wait()

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			log.Printf("Done server %s ...............\n", tb.ID())
			return nil
		case <-shutdownCh:
			log.Printf("Done server %s ...............\n", tb.ID())
			_ = tb.Stop()
			return nil
		case <-restartCh:
			log.Printf("Done server %s ...............\n", tb.ID())
			_ = tb.Stop()

			time.Sleep(5 * time.Second)
			reload()
			return nil
		case <-watchCh:
			log.Printf("Config change detected, restarting...\n")
			_ = tb.Stop()

			time.Sleep(5 * time.Second)
			reload()
			return nil
		}
	})
	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (tb *ToyBox) Stop() error {
	if tb.cancel != nil {
		tb.cancel()
	}
	return nil
}

func reload() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	log.Printf("execPath=%s args = %+v\n", execPath, os.Args)

	cmd := exec.Command(execPath, os.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start new process: %v", err)
	}
	log.Printf("Starting new process with PID: %d\n", cmd.Process.Pid)

	return nil
}
