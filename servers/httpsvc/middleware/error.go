package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GinGlobalException() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch e := err.(type) {
				default:
					c.JSON(500, gin.H{"code": 500, "msg": fmt.Sprintf("%v", e)})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}
