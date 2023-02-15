package util

import (
	"encoding/json"

	"github.com/google/uuid"
)

func JsonMarshal[T interface{}](str string) (T, bool) {
	var res T
	err := json.Unmarshal([]byte(str), &res)
	if err != nil {
		return res, false
	}
	return res, true
}

func Max[T int | int32 | int64 | float32 | float64 | uint | uint32 | uint64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T int | int32 | int64 | float32 | float64 | uint | uint32 | uint64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Add[T int | int32 | int64 | float32 | float64 | uint | uint32 | uint64 | string](a, b T) T {
	return a + b
}

func UUID(n ...int) string {
	if len(n) > 0 {
		return uuid.New().String()[:n[0]]
	}
	return uuid.New().String()[:8]
}
