package toybox

import "context"

type Server interface {
	toyboxBase
	Run(ctx context.Context) error
}
