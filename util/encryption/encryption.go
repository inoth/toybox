package encryption

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"os"
)

// https://segmentfault.com/a/1190000017346458
// 加解密算法：对称性加密算法、非对称性加密算法、散列算法，其中散列算法不可逆，无法解密，故而只能用于签名校验、身份验证
// 对称性加密算法：DES、3DES、AES
// 非对称性加密算法：RSA、DSA、ECC
// 散列算法：MD5、SHA1、HMAC
// GenerateRSAKey 生成私钥和公钥, bits参数指定证书大小
// 也可以直接通过openssl命令生成：
// 私钥：openssl genrsa -out rsa_private_key.pem 2048
// 公钥：openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
func GenerateRSAKey(bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	privateFile, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	defer func() {
		_ = privateFile.Close()
	}()
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	err = pem.Encode(privateFile, &privateBlock)
	if err != nil {
		return err
	}
	publicKey := privateKey.PublicKey
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err
	}
	publicFile, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	defer func() {
		_ = publicFile.Close()
	}()
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	err = pem.Encode(publicFile, &publicBlock)
	if err != nil {
		return err
	}
	return nil
}

// EncryptWithRSA rsa加密
func EncryptWithRSA(plainText string, publicKeyPath string) (string, error) {
	keyFile, err := os.Open(publicKeyPath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = keyFile.Close()
	}()
	info, _ := keyFile.Stat()
	buf := make([]byte, info.Size())
	_, err = keyFile.Read(buf)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(buf)
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	cipherBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(plainText))
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(cipherBytes), nil
}

// DecryptWithRSA rsa解密
func DecryptWithRSA(cipherText string, privateKeyPath string) (string, error) {
	keyFile, err := os.Open(privateKeyPath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = keyFile.Close()
	}()
	info, _ := keyFile.Stat()
	buf := make([]byte, info.Size())
	_, err = keyFile.Read(buf)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(buf)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	cipherBytes, err := base64.RawURLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	plainBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherBytes)
	if err != nil {
		return "", err
	}
	return string(plainBytes), nil
}

// EncryptWithSha256 sha256加密
func EncryptWithSha256(data string) string {
	tmp := base64.StdEncoding.EncodeToString([]byte(data))
	h := sha256.New()
	h.Write([]byte(tmp))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

func Md5(str string) (string, error) {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil)), nil
}

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
