package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	SIGNKEY = "BA5ktbKaV47uOcQpnuUT76GvBRYpMdHX"
)

type CustomerInfo struct {
	UserInfo map[string]interface{}
}

type CustomClaims struct {
	jwt.RegisteredClaims
	CustomerInfo
}

func CreateToken(userInfo map[string]interface{}, expire ...int64) (string, error) {
	key := []byte(SIGNKEY)
	c := CustomClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
			Issuer:    userInfo["name"].(string),
		},
		CustomerInfo{
			UserInfo: userInfo,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	sign, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return sign, nil
}

func ParseToken(tokenStr string) (*CustomerInfo, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SIGNKEY), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return &claims.CustomerInfo, nil
	} else {
		return nil, err
	}
}
