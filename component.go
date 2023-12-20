package toybox

import "context"

const (
	ComponentStatusOK      = "ready"
	ComponentStatusInvalid = "invalid"
)

type Component interface {
	Status() string
	Init(ctx context.Context) error
}
