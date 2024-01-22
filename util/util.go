package util

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
	letterBytes   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	src = rand.NewSource(time.Now().UnixNano())
)

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

func RandStr(ns ...int) string {
	n := First(16, ns)
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}

func getValueWithType(m map[string]interface{}, key string) (interface{}, bool) {
	val, ok := m[key]
	if !ok {
		return nil, false
	}
	return val, true
}

func GetIntValue(m map[string]interface{}, key string) (int, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return 0, false
	}
	if iVal, ok := val.(int); ok {
		return iVal, true
	}
	return 0, false
}

func GetFloatValue(m map[string]interface{}, key string) (float64, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return 0, false
	}
	if fVal, ok := val.(float64); ok {
		return fVal, true
	}
	return 0, false
}

func GetStringValue(m map[string]interface{}, key string) (string, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return "", false
	}
	if sVal, ok := val.(string); ok {
		return sVal, true
	}
	return "", false
}

func GetBoolValue(m map[string]interface{}, key string) (bool, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return false, false
	}
	if bVal, ok := val.(bool); ok {
		return bVal, true
	}
	return false, false
}

func GetStringSlice(m map[string]interface{}, key string) ([]string, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return nil, false
	}
	if ssVal, ok := val.([]string); ok {
		return ssVal, true
	}
	return nil, false
}

func GetInterfaceSlice(m map[string]interface{}, key string) ([]interface{}, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return nil, false
	}
	if ssVal, ok := val.([]interface{}); ok {
		return ssVal, true
	}
	return nil, false
}
