package toybox

import (
	"context"
	"fmt"
	"log"
	"os"
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

const defaultStopTimeout = 10 * time.Second

type ToyBox struct {
	option

	id     string
	ctx    context.Context
	cancel context.CancelFunc
}

func New(opts ...Option) *ToyBox {
	o := option{
		sigs:        []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		stopTimeout: defaultStopTimeout,
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

	// Signal channels
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, t.sigs...)
	defer signal.Stop(shutdownCh)

	eg, ctx := errgroup.WithContext(t.ctx)

	// Decode config into transports before starting.
	for _, tr := range t.transports {
		if t.cfg != nil {
			if matcher, ok := tr.(conf.ConfigureMatcher); ok {
				if err := t.cfg.PrimitiveDecode(matcher); err != nil {
					return errors.Wrap(err, "PrimitiveDecode")
				}
			}
		}
	}

	// Start all transports with panic recovery.
	for _, tr := range t.transports {
		eg.Go(func() error {
			return t.safeStart(ctx, tr)
		})
	}

	// Graceful stop: once context is cancelled, stop transports with timeout.
	// Registered before registration so Stop is always called on failure paths.
	eg.Go(func() error {
		<-ctx.Done()
		return t.stopAll()
	})

	// Signal handler goroutine — always active.
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case sig := <-shutdownCh:
			log.Printf("Received signal %v, shutting down server %s", sig, t.id)
		}
		t.cancel()
		return nil
	})

	// Service registration.
	svc := t.buildServiceInstance()
	if t.registrar != nil && svc != nil {
		if err := t.registrar.Register(t.ctx, svc); err != nil {
			t.cancel()
			_ = eg.Wait()
			return errors.Wrap(err, "service register")
		}
		log.Printf("Registered service %s/%s", svc.Name, svc.ID)
	}

	// Config watching (hot-reload).
	if watcher, ok := t.cfg.(interface {
		Watch(context.Context) error
	}); ok {
		eg.Go(func() error {
			return watcher.Watch(ctx)
		})
	}

	// Block until all goroutines finish.
	err := eg.Wait()

	// Deregister service after everything stopped.
	if t.registrar != nil && svc != nil {
		dCtx, dCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer dCancel()
		log.Printf("Deregistering service %s/%s", svc.Name, svc.ID)
		if dErr := t.registrar.Deregister(dCtx, svc); dErr != nil {
			log.Printf("Failed to deregister service: %v", dErr)
		}
	}

	log.Printf("Server %s stopped", t.id)
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

// safeStart runs transport.Start with panic recovery.
// Returns nil on context cancellation (expected shutdown), preserving real errors.
func (t *ToyBox) safeStart(ctx context.Context, tr interface{ Start(context.Context) error }) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("transport panic: %v", r)
			log.Printf("Recovered from transport panic: %v", r)
			t.cancel()
		}
	}()
	err = tr.Start(ctx)
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

// stopAll gracefully stops all transports with a timeout.
func (t *ToyBox) stopAll() error {
	stopCtx, stopCancel := context.WithTimeout(context.Background(), t.stopTimeout)
	defer stopCancel()

	var (
		mu   sync.Mutex
		errs []error
		wg   sync.WaitGroup
	)

	for _, tr := range t.transports {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					mu.Lock()
					errs = append(errs, fmt.Errorf("transport stop panic: %v", r))
					mu.Unlock()
				}
			}()
			if err := tr.Stop(stopCtx); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-stopCtx.Done():
		log.Printf("Transport stop timed out after %v", t.stopTimeout)
		return fmt.Errorf("transport stop timed out after %v", t.stopTimeout)
	}

	if len(errs) > 0 {
		return fmt.Errorf("transport stop errors: %v", errs)
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
