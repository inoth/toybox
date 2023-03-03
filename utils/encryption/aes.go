package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

//PKCS7 填充模式
func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//填充的反向操作，删除填充字符串
func pkcs7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("encryption str len 0 err")
	} else {
		unpadding := int(origData[length-1])
		return origData[:(length - unpadding)], nil
	}
}

//实现加密
func aesEcrypt(origData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	origData = pkcs7Padding(origData, blockSize)
	//采用AES加密方法中CBC加密模式
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))

	blocMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//实现解密
func aesDeCrypt(cypted []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块大小
	blockSize := block.BlockSize()
	//创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))
	//这个函数也可以用来解密
	blockMode.CryptBlocks(origData, cypted)
	//去除填充字符串
	origData, err = pkcs7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, err
}

//加密base64
func AesEcrypt(data []byte, key []byte) (string, error) {
	result, err := aesEcrypt(data, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

//解密
func AesDeCrypt(data string, key []byte) ([]byte, error) {
	//解密base64字符串
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	//执行AES解密
	return aesDeCrypt(dataByte, key)
}
