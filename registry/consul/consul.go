package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/inoth/toybox/registry"
)

// Option configures a consul registry.
type Option func(*Registry)

// WithHealthCheck enables a TTL-based health check with the given interval.
func WithHealthCheck(interval time.Duration) Option {
	return func(r *Registry) {
		r.healthCheckInterval = interval
		r.enableHealthCheck = true
	}
}

// Registry implements registry.Registrar and registry.Discovery backed by Consul.
type Registry struct {
	client              *api.Client
	enableHealthCheck   bool
	healthCheckInterval time.Duration
}

// New creates a new Consul registry.
func New(client *api.Client, opts ...Option) *Registry {
	r := &Registry{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Register registers the service instance with Consul.
func (r *Registry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	host, portStr, endpoint := "", "", ""
	if len(service.Endpoints) > 0 {
		endpoint = service.Endpoints[0]
	}
	if endpoint != "" {
		var err error
		host, portStr, err = net.SplitHostPort(endpoint)
		if err != nil {
			host = endpoint
		}
	}
	port, _ := strconv.Atoi(portStr)

	meta := service.Metadata
	if meta == nil {
		meta = make(map[string]string)
	}
	// Store full instance as JSON in a meta tag for richer deserialization.
	data, _ := json.Marshal(service)
	meta["__instance__"] = string(data)

	reg := &api.AgentServiceRegistration{
		ID:      service.ID,
		Name:    service.Name,
		Tags:    []string{service.Version},
		Port:    port,
		Address: host,
		Meta:    meta,
	}

	if r.enableHealthCheck {
		reg.Check = &api.AgentServiceCheck{
			TTL:                            r.healthCheckInterval.String(),
			DeregisterCriticalServiceAfter: (r.healthCheckInterval * 10).String(),
		}
	}

	if err := r.client.Agent().ServiceRegister(reg); err != nil {
		return fmt.Errorf("consul register: %w", err)
	}

	if r.enableHealthCheck {
		go r.keepAlive(ctx, service.ID)
	}
	return nil
}

func (r *Registry) keepAlive(ctx context.Context, serviceID string) {
	checkID := "service:" + serviceID
	ticker := time.NewTicker(r.healthCheckInterval / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = r.client.Agent().PassTTL(checkID, "")
		}
	}
}

// Deregister removes the service instance from Consul.
func (r *Registry) Deregister(_ context.Context, service *registry.ServiceInstance) error {
	return r.client.Agent().ServiceDeregister(service.ID)
}

// GetService returns all healthy instances for the given service name.
func (r *Registry) GetService(_ context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("consul get service: %w", err)
	}
	instances := make([]*registry.ServiceInstance, 0, len(entries))
	for _, entry := range entries {
		// Try rich deserialization from meta tag first.
		if raw, ok := entry.Service.Meta["__instance__"]; ok {
			var svc registry.ServiceInstance
			if json.Unmarshal([]byte(raw), &svc) == nil {
				instances = append(instances, &svc)
				continue
			}
		}
		// Fallback to reconstructing from Consul fields.
		svc := &registry.ServiceInstance{
			ID:       entry.Service.ID,
			Name:     entry.Service.Service,
			Metadata: entry.Service.Meta,
			Endpoints: []string{
				net.JoinHostPort(entry.Service.Address, strconv.Itoa(entry.Service.Port)),
			},
		}
		if len(entry.Service.Tags) > 0 {
			svc.Version = entry.Service.Tags[0]
		}
		instances = append(instances, svc)
	}
	return instances, nil
}

// Watch returns a watcher that monitors service instance changes via Consul blocking queries.
func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return &watcher{
		registry:    r,
		ctx:         ctx,
		serviceName: serviceName,
	}, nil
}

type watcher struct {
	registry    *Registry
	ctx         context.Context
	serviceName string
	lastIndex   uint64
	started     bool
}

// Next blocks until a service change occurs and returns the current list of instances.
func (w *watcher) Next() ([]*registry.ServiceInstance, error) {
	if !w.started {
		w.started = true
		instances, err := w.registry.GetService(w.ctx, w.serviceName)
		return instances, err
	}
	// Use Consul blocking query to wait for changes.
	opts := &api.QueryOptions{
		WaitIndex: w.lastIndex,
		WaitTime:  55 * time.Second,
	}
	opts = opts.WithContext(w.ctx)
	entries, meta, err := w.registry.client.Health().Service(w.serviceName, "", true, opts)
	if err != nil {
		return nil, fmt.Errorf("consul watch: %w", err)
	}
	w.lastIndex = meta.LastIndex

	instances := make([]*registry.ServiceInstance, 0, len(entries))
	for _, entry := range entries {
		if raw, ok := entry.Service.Meta["__instance__"]; ok {
			var svc registry.ServiceInstance
			if json.Unmarshal([]byte(raw), &svc) == nil {
				instances = append(instances, &svc)
				continue
			}
		}
		svc := &registry.ServiceInstance{
			ID:       entry.Service.ID,
			Name:     entry.Service.Service,
			Metadata: entry.Service.Meta,
			Endpoints: []string{
				net.JoinHostPort(entry.Service.Address, strconv.Itoa(entry.Service.Port)),
			},
		}
		if len(entry.Service.Tags) > 0 {
			svc.Version = entry.Service.Tags[0]
		}
		instances = append(instances, svc)
	}
	return instances, nil
}

// Stop is a no-op; stopping is driven by context cancellation.
func (w *watcher) Stop() error {
	return nil
}
