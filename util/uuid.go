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

func First[T any](defaultArg T, args []T) T {
	if len(args) > 0 {
		defaultArg = args[0]
	}
	return defaultArg
}

func UUID(ns ...int) string {
	n := First(16, ns)
	uuidStr := uuid.New().String()
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	return uuidStr[0:n]
}

func UUIDV7(ns ...int) string {
	n := First(16, ns)
	uuidV7, _ := uuid.NewV7()
	uuidStr := strings.ReplaceAll(uuidV7.String(), "-", "")
	return uuidStr[len(uuidStr)-n:]
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
