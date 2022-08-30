package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/inoth/ino-toybox/components/logger"
	"github.com/inoth/ino-toybox/res"
	"github.com/pkg/errors"
)

const SIGNKEY = "BA5ktbKaV47uOcQpnuUT76GvBRYpMdHX"

type CustomerInfo struct {
	Uid  string
	Name string
}
type CustomClaims struct {
	*jwt.StandardClaims
	*CustomerInfo
}

func CreateToken(uid string, name string, expire ...int64) (string, error) {
	key := []byte(SIGNKEY)
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
			Uid:  uid,
			Name: name,
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
		return []byte(SIGNKEY), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.CustomerInfo, nil
	} else {
		return nil, errors.Wrap(err, "")
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		token = c.GetHeader("Authorization")
		if token == "" {
			token, _ = c.Cookie("Authorization")
		}
		if token == "" {
			res.Unauthrized(c, "Unauthrized.")
			c.Abort()
			return
		}
		user, err := ParseToken(token)
		if err != nil {
			logger.Zap.Error(fmt.Sprintf("jwt解析失败：%v", err))
			logger.Zap.Error(fmt.Sprintf("无效token: %v", token))
			res.Unauthrized(c, "Unauthrized.")
			c.Abort()
			return
		}
		c.Set("USER_ID", user.Uid)
		c.Set("USER_NAME", user.Name)
		c.Next()
	}
}
