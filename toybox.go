package toybox

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/inoth/toybox/conf"
	"github.com/inoth/toybox/util"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type ToyBox struct {
	option

	id     string
	ctx    context.Context
	cancel context.CancelFunc
}

func New(opts ...Option) *ToyBox {
	o := option{
		sigs: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
	for _, opt := range opts {
		opt(&o)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &ToyBox{
		option: o,
		id:     util.UUID(),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (t *ToyBox) Run() error {
	log.Printf("Starting server ID:%s (PID: %d)\n", t.id, os.Getpid())

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, t.sigs...)

	restartCh := make(chan os.Signal, 1)
	signal.Notify(restartCh, syscall.SIGHUP)

	wg := sync.WaitGroup{}
	eg, ctx := errgroup.WithContext(t.ctx)

	for _, transport := range t.transports {
		if t.cfg != nil {
			if err := t.cfg.PrimitiveDecode(transport.(conf.ConfigureMatcher)); err != nil {
				return errors.Wrap(err, "PrimitiveDecode")
			}
		}
		eg.Go(func() error {
			<-ctx.Done()
			return transport.Stop(ctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			defer wg.Done()
			return transport.Start(ctx)
		})
	}
	wg.Wait()

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			log.Printf("Done server %s", t.id)
			return nil
		case <-shutdownCh:
			log.Printf("Shutting down server %s", t.id)
			if t.cancel != nil {
				t.cancel()
			}
		case <-restartCh:
			log.Printf("Restarting server %s", t.id)
			if t.cancel != nil {
				t.cancel()
			}
			time.Sleep(time.Second * 3)
			if err := reload(); err != nil {
				log.Printf("Failed to reload: %v", err)
			}
		}
		return nil
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
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
