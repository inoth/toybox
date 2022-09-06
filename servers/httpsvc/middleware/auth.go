package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/inoth/ino-toybox/components/logger"
	"github.com/inoth/ino-toybox/res"
	"github.com/inoth/ino-toybox/util/auth"
)

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
		user, err := auth.ParseToken(token)
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
