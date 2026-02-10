module github.com/inoth/toybox

go 1.25

retract (
	[v1.0.0, v1.1.9]
	[v0.0.1, v0.9.9]
)

require (
	github.com/google/uuid v1.6.0
	github.com/pkg/errors v0.9.1
	golang.org/x/sync v0.19.0
)
