package toybox

import "context"

type Component interface {
	toyboxBase
	Init(ctx context.Context) error
}
