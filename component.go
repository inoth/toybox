package toybox

import "context"

type Component interface {
	base
	Init(ctx context.Context) error
}
