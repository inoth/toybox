package toybox

import "context"

type Component interface {
	basic
	Init(ctx context.Context) error
}
