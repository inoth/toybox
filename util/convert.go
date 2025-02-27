package util

import (
	"fmt"
	"unsafe"
)

func ConvertWithUnsafe[I, O any](input *I) (output *O, err error) {
	defer func() {
		if err := recover(); err != nil {
			err = fmt.Errorf("%v", err)
		}
	}()
	if input == nil {
		err = fmt.Errorf("input object is empty")
		return
	}
	output = (*O)(unsafe.Pointer(input))
	return
}
