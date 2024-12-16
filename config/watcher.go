package config

import "context"

type Watcher interface {
	Next(ctx context.Context)
	Probe(chan<- struct{})
}
