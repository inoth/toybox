package toybox

import "context"

type Server interface {
	Ready() bool
	Run(ctx context.Context) error
}
