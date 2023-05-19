package component

type Component interface {
	Name() string
	Init() error
	String() string
}
