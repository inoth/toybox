package convert

import (
	"fmt"
	"unsafe"
)

// 相同接口结构体之间转换
func ConvertByUnsafe[A, B any](input *A) (output *B, err error) {
	if input == nil {
		err = fmt.Errorf("input struct is empty")
		return
	}
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	output = (*B)(unsafe.Pointer(input))
	return
}

// 结构体必须结构和类型一致
func ConvertWithUnsafe[A, B any](src *A) (dst B, ok bool) {
	if src == nil {
		return
	}
	dst = *(*B)(unsafe.Pointer(src))
	return
}
