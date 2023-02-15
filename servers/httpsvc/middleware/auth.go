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
			token = c.Query("token")
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
		c.Set("USER_ID", user.UserInfo["uid"])
		c.Set("USER_NAME", user.UserInfo["name"])
		c.Set("USER_AVATER", user.UserInfo["avatar"])
		c.Next()
	}
}
