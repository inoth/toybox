package common

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5WithString(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Md5(buf []byte) string {
	h := md5.New()
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}
