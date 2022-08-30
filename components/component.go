package components

/*
	常用组件, 比如mysql、redis、config等...
*/

type Component interface {
	Init() error
}
