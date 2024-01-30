package toybox

import "context"

type Server interface {
	basic
	Run(ctx context.Context) error
}
