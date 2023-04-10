package test

import (
	"fmt"
	"testing"

	"github.com/inoth/toybox/utils"
)

func TestSubStr(t *testing.T) {
	a := "123456789"
	fmt.Println(utils.Substr(a, 0, 5))
}
