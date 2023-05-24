package middleware

import (
	"fmt"
	"time"

	"github.com/inoth/toybox/common"

	"github.com/golang-jwt/jwt/v4"
)

type CustomerInfo struct {
	UserInfo map[string]interface{}
}

type CustomClaims struct {
	jwt.RegisteredClaims
	CustomerInfo
}

func CreateToken(signKey string, userInfo map[string]interface{}, expire ...time.Duration) (string, error) {
	key := []byte(signKey)
	exp := time.Duration(24)
	if len(expire) > 0 {
		exp = expire[0]
	}
	issuser, _ := common.GetStringValue(userInfo, "name")
	c := CustomClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * exp)),
			Issuer:    issuser,
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

func ParseToken(signKey, tokenStr string) (*CustomerInfo, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signKey), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return &claims.CustomerInfo, nil
	} else {
		return nil, err
	}
}
