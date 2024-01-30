package toybox

type basic interface {
	Name() string
	Ready() bool
	IsReady()
}
