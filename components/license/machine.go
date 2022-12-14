package license

import (
	"crypto/md5"
	"encoding/hex"
)

func MachineCode() (string, error) {
	unique, err := GetUnique()
	if err != nil {
		return "", err
	}
	h := md5.New()
	h.Write([]byte(unique))
	return hex.EncodeToString(h.Sum(nil)), nil
}
