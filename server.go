package toybox

import "context"

type Server interface {
	Name() string
	Ready() bool
	IsReady()
	Run(ctx context.Context) error
}
