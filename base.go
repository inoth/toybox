package toybox

type base interface {
	Name() string
	Ready() bool
	IsReady()
}
