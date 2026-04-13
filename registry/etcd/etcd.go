package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/inoth/toybox/registry"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const defaultPrefix = "/toybox/services"

// Option configures an etcd registry.
type Option func(*Registry)

// WithPrefix sets the key prefix in etcd.
func WithPrefix(prefix string) Option {
	return func(r *Registry) {
		r.prefix = prefix
	}
}

// WithLeaseTTL sets the lease TTL in seconds. Default is 15.
func WithLeaseTTL(ttl int64) Option {
	return func(r *Registry) {
		r.ttl = ttl
	}
}

// Registry implements registry.Registrar and registry.Discovery backed by etcd.
type Registry struct {
	client *clientv3.Client
	prefix string
	ttl    int64
	lease  clientv3.LeaseID
}

// New creates a new etcd registry.
func New(client *clientv3.Client, opts ...Option) *Registry {
	r := &Registry{
		client: client,
		prefix: defaultPrefix,
		ttl:    15,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Registry) serviceKey(service *registry.ServiceInstance) string {
	return fmt.Sprintf("%s/%s/%s", r.prefix, service.Name, service.ID)
}

// Register registers the service instance with a keep-alive lease.
func (r *Registry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	data, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("marshal service: %w", err)
	}
	grant, err := r.client.Grant(ctx, r.ttl)
	if err != nil {
		return fmt.Errorf("etcd grant lease: %w", err)
	}
	r.lease = grant.ID

	key := r.serviceKey(service)
	_, err = r.client.Put(ctx, key, string(data), clientv3.WithLease(grant.ID))
	if err != nil {
		return fmt.Errorf("etcd put: %w", err)
	}

	ch, err := r.client.KeepAlive(ctx, grant.ID)
	if err != nil {
		return fmt.Errorf("etcd keepalive: %w", err)
	}
	// Drain keepalive responses in the background.
	go func() {
		for range ch {
		}
	}()
	return nil
}

// Deregister removes the service instance from etcd.
func (r *Registry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	key := r.serviceKey(service)
	_, err := r.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("etcd delete: %w", err)
	}
	if r.lease != 0 {
		r.client.Revoke(ctx, r.lease)
	}
	return nil
}

// GetService returns all instances for the given service name.
func (r *Registry) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	prefix := fmt.Sprintf("%s/%s/", r.prefix, serviceName)
	resp, err := r.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("etcd get: %w", err)
	}
	instances := make([]*registry.ServiceInstance, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var svc registry.ServiceInstance
		if err := json.Unmarshal(kv.Value, &svc); err != nil {
			continue
		}
		instances = append(instances, &svc)
	}
	return instances, nil
}

// Watch returns a watcher that monitors service instance changes.
func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	prefix := fmt.Sprintf("%s/%s/", r.prefix, serviceName)
	w := &watcher{
		client:      r.client,
		prefix:      prefix,
		ctx:         ctx,
		serviceName: serviceName,
		registry:    r,
	}
	return w, nil
}

type watcher struct {
	client      *clientv3.Client
	prefix      string
	ctx         context.Context
	serviceName string
	registry    *Registry
	watchCh     clientv3.WatchChan
	started     bool
}

// Next blocks until a service change occurs and returns the current list of instances.
func (w *watcher) Next() ([]*registry.ServiceInstance, error) {
	if !w.started {
		w.started = true
		// Return current state on first call.
		return w.registry.GetService(w.ctx, w.serviceName)
	}
	if w.watchCh == nil {
		w.watchCh = w.client.Watch(w.ctx, w.prefix, clientv3.WithPrefix())
	}
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case resp, ok := <-w.watchCh:
		if !ok {
			// Channel closed, retry after brief pause.
			time.Sleep(time.Second)
			w.watchCh = nil
			return w.registry.GetService(w.ctx, w.serviceName)
		}
		if resp.Err() != nil {
			return nil, resp.Err()
		}
		return w.registry.GetService(w.ctx, w.serviceName)
	}
}

// Stop cancels the watcher. The parent context cancellation also stops it.
func (w *watcher) Stop() error {
	return nil
}
