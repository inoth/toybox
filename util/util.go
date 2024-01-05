package util

import "github.com/google/uuid"

func Max[T int | int32 | int64 | float32 | float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T int | int32 | int64 | float32 | float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func First[T interface{}](defaultArg T, args []T) T {
	if len(args) > 0 {
		defaultArg = args[0]
	}
	return defaultArg
}

func UUID() string {
	return uuid.New().String()
}
