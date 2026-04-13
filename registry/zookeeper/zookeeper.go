package zookeeper

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/inoth/toybox/registry"
)

const defaultBasePath = "/toybox/services"

// Option configures a zookeeper registry.
type Option func(*Registry)

// WithBasePath sets the base znode path.
func WithBasePath(base string) Option {
	return func(r *Registry) {
		r.basePath = base
	}
}

// Registry implements registry.Registrar and registry.Discovery backed by ZooKeeper.
type Registry struct {
	conn     *zk.Conn
	basePath string
}

// New creates a new ZooKeeper registry.
func New(conn *zk.Conn, opts ...Option) *Registry {
	r := &Registry{
		conn:     conn,
		basePath: defaultBasePath,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// ensurePath creates the znode path recursively if it does not exist.
func (r *Registry) ensurePath(p string) error {
	if exists, _, err := r.conn.Exists(p); err != nil {
		return err
	} else if exists {
		return nil
	}
	parent := path.Dir(p)
	if parent != "/" && parent != "." {
		if err := r.ensurePath(parent); err != nil {
			return err
		}
	}
	_, err := r.conn.Create(p, nil, 0, zk.WorldACL(zk.PermAll))
	if err == zk.ErrNodeExists {
		return nil
	}
	return err
}

func (r *Registry) servicePath(serviceName string) string {
	return path.Join(r.basePath, serviceName)
}

func (r *Registry) instancePath(service *registry.ServiceInstance) string {
	return path.Join(r.basePath, service.Name, service.ID)
}

// Register creates an ephemeral znode for the service instance.
func (r *Registry) Register(_ context.Context, service *registry.ServiceInstance) error {
	svcPath := r.servicePath(service.Name)
	if err := r.ensurePath(svcPath); err != nil {
		return fmt.Errorf("zk ensure path: %w", err)
	}
	data, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("marshal service: %w", err)
	}
	nodePath := r.instancePath(service)
	_, err = r.conn.Create(nodePath, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err == zk.ErrNodeExists {
		// Update existing node.
		_, stat, _ := r.conn.Get(nodePath)
		_, err = r.conn.Set(nodePath, data, stat.Version)
	}
	if err != nil {
		return fmt.Errorf("zk register: %w", err)
	}
	return nil
}

// Deregister removes the service instance znode.
func (r *Registry) Deregister(_ context.Context, service *registry.ServiceInstance) error {
	nodePath := r.instancePath(service)
	_, stat, err := r.conn.Get(nodePath)
	if err != nil {
		if err == zk.ErrNoNode {
			return nil
		}
		return fmt.Errorf("zk get: %w", err)
	}
	return r.conn.Delete(nodePath, stat.Version)
}

// GetService returns all registered instances for the given service name.
func (r *Registry) GetService(_ context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	svcPath := r.servicePath(serviceName)
	children, _, err := r.conn.Children(svcPath)
	if err != nil {
		if err == zk.ErrNoNode {
			return nil, nil
		}
		return nil, fmt.Errorf("zk children: %w", err)
	}
	instances := make([]*registry.ServiceInstance, 0, len(children))
	for _, child := range children {
		data, _, err := r.conn.Get(path.Join(svcPath, child))
		if err != nil {
			continue
		}
		var svc registry.ServiceInstance
		if json.Unmarshal(data, &svc) == nil {
			instances = append(instances, &svc)
		}
	}
	return instances, nil
}

// Watch returns a watcher that monitors service instance changes via ZooKeeper watches.
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
	started     bool
}

// Next blocks until a service change occurs and returns the current list of instances.
func (w *watcher) Next() ([]*registry.ServiceInstance, error) {
	if !w.started {
		w.started = true
		return w.registry.GetService(w.ctx, w.serviceName)
	}
	svcPath := w.registry.servicePath(w.serviceName)
	for {
		_, _, eventCh, err := w.registry.conn.ChildrenW(svcPath)
		if err != nil {
			if err == zk.ErrNoNode {
				time.Sleep(time.Second)
				continue
			}
			return nil, fmt.Errorf("zk watch: %w", err)
		}
		select {
		case <-w.ctx.Done():
			return nil, w.ctx.Err()
		case <-eventCh:
			return w.registry.GetService(w.ctx, w.serviceName)
		}
	}
}

// Stop is a no-op; stopping is driven by context cancellation.
func (w *watcher) Stop() error {
	return nil
}
