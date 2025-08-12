package convert

import (
	"unsafe"
)

// 相同接口结构体之间转换
func ConvertByUnsafe[A, B any](input *A) (output *B, ok bool) {
	if input == nil {
		return
	}
	output = (*B)(unsafe.Pointer(input))
	return
}
