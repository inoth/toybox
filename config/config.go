package config

import "context"

const (
	DefaultDir = "config"
)

type ConfigMate interface {
	PrimitiveDecode(vals ...ConfigureMatcher) error
}

type ConfigureMatcher interface {
	Name() string
}

type Source interface {
	Load(format string) (string, error)
}

type Watcher interface {
	Next(ctx context.Context)
	Probe(chan<- struct{})
}
