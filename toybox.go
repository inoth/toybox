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
const defaultPIDFile = "toybox.pid"

// errReload is a sentinel error that signals a transport reload cycle.
var errReload = errors.New("reload")

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
		pidFile:     defaultPIDFile,
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
	// Write PID file if configured (nginx-style, enables: kill -HUP $(cat pidfile)).
	if t.pidFile != "" {
		if err := t.writePIDFile(); err != nil {
			return errors.Wrap(err, "write pid file")
		}
		defer t.removePIDFile()
	}

	log.Printf("Starting server ID:%s (PID: %d)\n", t.id, os.Getpid())
	defer t.cancel() // ensure background goroutines are cleaned up on exit

	// Shutdown signals.
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, t.sigs...)
	defer signal.Stop(shutdownCh)

	// Reload signal (SIGHUP, nginx-style).
	reloadCh := make(chan os.Signal, 1)
	signal.Notify(reloadCh, syscall.SIGHUP)
	defer signal.Stop(reloadCh)

	// Config change notification channel.
	configChangeCh := make(chan struct{}, 1)
	if t.cfg != nil {
		if onChanger, ok := t.cfg.(interface{ OnChange(func()) }); ok {
			onChanger.OnChange(func() {
				select {
				case configChangeCh <- struct{}{}:
				default:
				}
			})
		}
	}

	// Start config watcher in background (runs for the lifetime of the process).
	if watcher, ok := t.cfg.(interface {
		Watch(context.Context) error
	}); ok {
		go func() {
			if err := watcher.Watch(t.ctx); err != nil {
				log.Printf("config watch error: %v", err)
			}
		}()
	}

	// Service registration (once, for the process lifetime).
	svc := t.buildServiceInstance()
	if t.registrar != nil && svc != nil {
		if err := t.registrar.Register(t.ctx, svc); err != nil {
			return errors.Wrap(err, "service register")
		}
		log.Printf("Registered service %s/%s", svc.Name, svc.ID)
	}

	// Main lifecycle loop — each iteration is one transport generation.
	// On SIGHUP or config change, transports are gracefully stopped and restarted.
	var finalErr error
	for {
		err := t.runTransportCycle(shutdownCh, reloadCh, configChangeCh)
		if errors.Is(err, errReload) {
			log.Printf("Reloading server %s ...", t.id)
			// Reload config from source (idempotent if already reloaded by watcher).
			if reloader, ok := t.cfg.(interface{ Reload() error }); ok {
				if rErr := reloader.Reload(); rErr != nil {
					log.Printf("config reload error: %v", rErr)
				}
			}
			continue
		}
		finalErr = err
		break
	}

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
	if finalErr != nil && !errors.Is(finalErr, context.Canceled) {
		return finalErr
	}
	return nil
}

// runTransportCycle runs one generation of transports. It blocks until all
// transports finish. Returns errReload to signal that a new cycle should start,
// or any other error/nil for final shutdown.
func (t *ToyBox) runTransportCycle(shutdownCh, reloadCh <-chan os.Signal, configChangeCh <-chan struct{}) error {
	// Drain stale signals from a previous cycle.
	drainSignals(reloadCh, configChangeCh)

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

	cycleCtx, cycleCancel := context.WithCancel(t.ctx)
	defer cycleCancel()

	eg, ctx := errgroup.WithContext(cycleCtx)

	// Start all transports with panic recovery.
	for _, tr := range t.transports {
		eg.Go(func() error {
			return t.safeStart(ctx, tr)
		})
	}

	// Graceful stop: once this cycle's context is cancelled, drain all
	// in-flight work by calling Stop on every transport before returning.
	eg.Go(func() error {
		<-ctx.Done()
		return t.stopAll()
	})

	// Track whether this cycle ends due to a reload request.
	reloadRequested := make(chan struct{})

	// Signal / reload handler goroutine.
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case sig := <-shutdownCh:
			log.Printf("Received signal %v, shutting down server %s", sig, t.id)
			t.cancel() // cancel root context → full shutdown
			return nil
		case <-reloadCh:
			log.Printf("Received SIGHUP, reloading server %s", t.id)
			close(reloadRequested)
			cycleCancel() // cancel only this cycle → triggers graceful stop then restart
			return nil
		case <-configChangeCh:
			log.Printf("Config changed, reloading server %s", t.id)
			close(reloadRequested)
			cycleCancel()
			return nil
		}
	})

	err := eg.Wait()

	// Determine if this was a reload or a shutdown.
	select {
	case <-reloadRequested:
		return errReload
	default:
	}
	return err
}

// drainSignals discards any pending signals left over from a previous cycle.
func drainSignals(reloadCh <-chan os.Signal, configChangeCh <-chan struct{}) {
	for {
		select {
		case <-reloadCh:
		case <-configChangeCh:
		default:
			return
		}
	}
}

// writePIDFile writes the current process ID to the configured PID file.
func (t *ToyBox) writePIDFile() error {
	return os.WriteFile(t.pidFile, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0644)
}

// removePIDFile removes the PID file on shutdown.
func (t *ToyBox) removePIDFile() {
	if err := os.Remove(t.pidFile); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to remove pid file: %v", err)
	}
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
