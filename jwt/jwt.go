package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/inoth/toybox/util"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Uid    string                 `json:"uid"`
	KeyMap map[string]interface{} `json:"key_map,omitempty"`
}

func NewClaims(uid string, keyMap map[string]interface{}, expire ...time.Duration) *CustomClaims {
	exp := util.First(24, expire)
	return &CustomClaims{
		Uid:    uid,
		KeyMap: keyMap,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    uid,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * exp)),
		},
	}
}

func GenToken(cu *CustomClaims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cu)
	sign, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return sign, nil
}

func ParseToken(jwtStr, key string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(jwtStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
