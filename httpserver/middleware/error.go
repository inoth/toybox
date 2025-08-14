package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/httpserver/res"
	"github.com/inoth/toybox/logger"
)

func Recover(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := recover(); err != nil {
			switch e := err.(type) {
			case error:
				log.LogError(c, e.Error())
				res.Failed(c, e.Error())
			default:
				log.LogError(c, "InternalServerError")
				res.Failed(c, "InternalServerError")
			}
			c.Abort()
		}
	}
}
