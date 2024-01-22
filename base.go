package toybox

type toyboxBase interface {
	Name() string
	Ready() bool
	IsReady()
}
