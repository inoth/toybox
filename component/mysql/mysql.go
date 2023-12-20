package mysql

import "github/inoth/toybox"

type MysqlComponent struct {
}

func New(tb *toybox.ToyBox) *MysqlComponent {
	return &MysqlComponent{}
}
