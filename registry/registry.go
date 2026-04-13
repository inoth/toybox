package registry

import "context"

// ServiceInstance represents a registered service instance.
type ServiceInstance struct {
	ID        string
	Name      string
	Version   string
	Metadata  map[string]string
	Endpoints []string
}

// Registrar registers and deregisters service instances with a service registry.
type Registrar interface {
	Register(ctx context.Context, service *ServiceInstance) error
	Deregister(ctx context.Context, service *ServiceInstance) error
}

// Discovery discovers service instances from the registry.
type Discovery interface {
	GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

// Watcher watches for service instance changes.
type Watcher interface {
	Next() ([]*ServiceInstance, error)
	Stop() error
}

// Endpointer is an optional interface that transports can implement
// to expose their network endpoint for service registration.
type Endpointer interface {
	Endpoint() (string, error)
}
