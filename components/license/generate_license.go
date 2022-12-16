package license

import (
	"crypto"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/inoth/ino-toybox/util"
	"github.com/inoth/ino-toybox/util/encryption"
)

var ErrLicCode = errors.New("未知订阅密钥")

// 生成订阅码
func MakeLicense(token string, priKey []byte, license License) {
	str := license.String()
	if len(str) <= 0 {
		return
	}

	key, err := encryption.Md5(token)
	if err != nil {
		fmt.Printf("ERR: %v", err.Error())
		return
	}
	// 混淆码
	randStr := util.RandString(32)

	encrypt, err := encryption.AesEcrypt(str, []byte(key))
	if err != nil {
		fmt.Printf("ERR: %v", err.Error())
		return
	}

	sign, err := encryption.SignWithRSA([]byte(encrypt), priKey, crypto.SHA256)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	f, err := os.Create(key + ".license")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()

	var sb strings.Builder
	// sb.Grow(len(randStr) + len(encrypt))
	sb.WriteString(randStr)
	sb.WriteString("\\")
	sb.WriteString(encrypt)
	sb.WriteString("\\")
	sb.WriteString(sign)
	f.WriteString(sb.String())
}

// 解析订阅码
func ParseLicense(token, pubKey, lic string) (*License, error) {
	licSpli := strings.Split(lic, "\\")
	if len(licSpli) != 3 {
		return nil, ErrLicCode
	}
	encryptData := licSpli[1]
	sign := licSpli[2]

	data, err := encryption.AesDeCrypt(encryptData, []byte(token))
	if err != nil {
		return nil, ErrLicCode
	}
	if !encryption.VerifyWithRSA(data, []byte(pubKey), sign, crypto.SHA256) {
		return nil, ErrLicCode
	}
	var license License
	err = json.Unmarshal(data, &license)
	if err != nil {
		return nil, ErrLicCode
	}
	// 无法获取机器码
	machineCode, err := MachineCode()
	if err != nil {
		return nil, ErrLicCode
	}
	// 未授权机器
	if _, ok := license.LegalMachine[machineCode]; !ok {
		return nil, ErrLicCode
	}
	return &license, nil
}
