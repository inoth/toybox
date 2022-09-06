package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/inoth/ino-toybox/util"
	"github.com/pkg/errors"
)

type CustomerInfo struct {
	Uid    string
	Name   string
	Avater string
}
type CustomClaims struct {
	*jwt.StandardClaims
	*CustomerInfo
}

func CreateToken(uid, name, avater string, expire ...int64) (string, error) {
	key := []byte(util.SIGNKEY)
	expiresAt := time.Now().Add(time.Hour * 24).Unix()
	if len(expire) > 0 {
		expiresAt = expire[0]
	}
	c := CustomClaims{
		&jwt.StandardClaims{
			ExpiresAt: expiresAt,
			Issuer:    name,
		},
		&CustomerInfo{
			Uid:    uid,
			Name:   name,
			Avater: avater,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	sign, err := token.SignedString(key)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	return sign, nil
}

func ParseToken(tokenStr string) (*CustomerInfo, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v\n", token.Header["alg"])
		}
		return []byte(util.SIGNKEY), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.CustomerInfo, nil
	} else {
		return nil, errors.Wrap(err, "")
	}
}
