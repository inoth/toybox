package license

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/inoth/ino-toybox/util"
	"github.com/inoth/ino-toybox/util/encryption"
)

func makeLicense(machineCode string, license License) {
	l, err := json.Marshal(license)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	encrypt, err := encryption.AesEncrypt(l, []byte(machineCode))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	f, err := os.Create(machineCode + ".license")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()
	f.Write(encrypt)
}

func MakeLicense(aesKey string, license License) {
	str := license.String()
	if len(str) <= 0 {
		return
	}

	key, err := encryption.Md5(aesKey)
	if err != nil {
		fmt.Printf("ERR: %v", err.Error())
		return
	}
	// 混淆码
	randStr := util.RandString(32)

	encrypt, err := encryption.AesEncrypt(str, []byte(key))
	if err != nil {
		fmt.Printf("ERR: %v", err.Error())
		return
	}

	f, err := os.Create(key + ".license")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()

	var sb strings.Builder
	sb.Grow(len(randStr) + len(encrypt))
	sb.WriteString(randStr)
	sb.Write(encrypt)

	f.WriteString(sb.String())
}
