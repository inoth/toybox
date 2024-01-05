package toybox

import "context"

type Component interface {
	Name() string
	Ready() bool
	IsReady()
	Init(ctx context.Context) error
}
