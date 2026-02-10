package conf

import "context"

type ConfigMate interface {
	PrimitiveDecode(vals ...ConfigureMatcher) error
}

type ConfigureMatcher interface {
	TransportName() string
}

type Source interface {
	Load(wildcard string) (string, error)
}

type Watcher interface {
	Watche(ctx context.Context, trigger chan<- struct{})
}
