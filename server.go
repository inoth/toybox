package toybox

import "context"

type Server interface {
	base
	Run(ctx context.Context) error
}
