package encryption

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

// rsa公钥加密
func EncryptWithRSAString(plainText, publicKey string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	cipherBytes, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(plainText))
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(cipherBytes), nil
}

// rsa私钥解密
func DecryptWithRSAString(cipherText, privateKey string) (string, error) {
	block, _ := pem.Decode([]byte(privateKey))
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	cipherBytes, err := base64.RawURLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	plainBytes, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, cipherBytes)
	if err != nil {
		return "", err
	}
	return string(plainBytes), nil
}

// 私钥签名
func SignWithRSA(data, privateKey []byte, sHash crypto.Hash) (string, error) {
	block, _ := pem.Decode(privateKey)
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	hash := sHash.New()
	hash.Write(data)
	sign, err := rsa.SignPKCS1v15(rand.Reader, priKey, sHash, hash.Sum(nil))
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(sign), nil
}

// 公钥验签
func VerifyWithRSA(data, publicKey []byte, sign string, sHash crypto.Hash) bool {
	block, _ := pem.Decode(publicKey)
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return false
	}
	h := sHash.New()
	h.Write(data)
	orignSign, _ := base64.RawStdEncoding.DecodeString(sign)
	err = rsa.VerifyPKCS1v15(pubKey, sHash, h.Sum(nil), orignSign)
	return err == nil
}
