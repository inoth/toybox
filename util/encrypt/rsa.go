package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/inoth/toybox/util"
)

func parsePrivateKey(privateKeyStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return nil, fmt.Errorf("failed to parse private key")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

func parsePublicKey(publicKeyStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return nil, fmt.Errorf("failed to parse public key")
	}
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

func encryptRSA(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, err
	}
	return encryptedData, nil
}

func decryptRSA(privateKey *rsa.PrivateKey, encryptedData []byte) ([]byte, error) {
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedData)
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}

func EncryptRSA(publicKey string, data []byte) ([]byte, error) {
	pubKey, err := parsePublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return encryptRSA(pubKey, data)
}

func DecryptRSA(privateKey string, encryptedData []byte) ([]byte, error) {
	privKey, err := parsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return decryptRSA(privKey, encryptedData)
}

func EncryptRSAFile(pubKeyPath string, data []byte) ([]byte, error) {
	publicKey, err := util.ReadFile(pubKeyPath)
	if err != nil {
		return nil, err
	}
	pubKey, err := parsePublicKey(string(publicKey))
	if err != nil {
		return nil, err
	}
	return encryptRSA(pubKey, data)
}

func DecryptRSAFile(privKeyPath string, encryptedData []byte) ([]byte, error) {
	privateKey, err := util.ReadFile(privKeyPath)
	if err != nil {
		return nil, err
	}
	privKey, err := parsePrivateKey(string(privateKey))
	if err != nil {
		return nil, err
	}
	return decryptRSA(privKey, encryptedData)
}
