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

	"github.com/inoth/toybox/bootstrap"
	"github.com/inoth/toybox/conf"
	"github.com/inoth/toybox/registry"
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
	// Apply bootstrap config for fields not explicitly set.
	if o.bootstrap != nil {
		applyBootstrap(&o, o.bootstrap)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &ToyBox{
		option: o,
		id:     util.UUID(),
		ctx:    ctx,
		cancel: cancel,
	}
}

// applyBootstrap fills in option fields from bootstrap config when not already set.
func applyBootstrap(o *option, cfg *bootstrap.Config) {
	if o.serviceName == "" && cfg.ServiceName != "" {
		o.serviceName = cfg.ServiceName
	}
	if o.serviceVersion == "" && cfg.ServiceVersion != "" {
		o.serviceVersion = cfg.ServiceVersion
	}
	if o.registrar == nil {
		r, err := bootstrap.SetupRegistry(cfg)
		if err != nil {
			log.Printf("bootstrap: setup registry failed: %v", err)
		} else if r != nil {
			o.registrar = r
		}
	}
	if o.cfg == nil {
		c, err := bootstrap.SetupConfig(cfg)
		if err != nil {
			log.Printf("bootstrap: setup config failed: %v", err)
		} else if c != nil {
			o.cfg = c
		}
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

	// Service registration
	svc := t.buildServiceInstance()
	if t.registrar != nil && svc != nil {
		if err := t.registrar.Register(t.ctx, svc); err != nil {
			return errors.Wrap(err, "service register")
		}
		log.Printf("Registered service %s/%s", svc.Name, svc.ID)
		eg.Go(func() error {
			<-ctx.Done()
			dCtx, dCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer dCancel()
			log.Printf("Deregistering service %s/%s", svc.Name, svc.ID)
			return t.registrar.Deregister(dCtx, svc)
		})
	}

	// Config watching (hot-reload)
	if watcher, ok := t.cfg.(interface {
		Watch(context.Context) error
	}); ok {
		eg.Go(func() error {
			return watcher.Watch(ctx)
		})
	}

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

func (t *ToyBox) buildServiceInstance() *registry.ServiceInstance {
	if t.serviceName == "" {
		return nil
	}
	endpoints := make([]string, 0, len(t.transports))
	for _, tr := range t.transports {
		if e, ok := tr.(registry.Endpointer); ok {
			ep, err := e.Endpoint()
			if err == nil {
				endpoints = append(endpoints, ep)
			}
		}
	}
	return &registry.ServiceInstance{
		ID:        t.id,
		Name:      t.serviceName,
		Version:   t.serviceVersion,
		Metadata:  t.metadata,
		Endpoints: endpoints,
	}
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
