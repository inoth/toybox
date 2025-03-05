package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/ginserver/res"
	"github.com/inoth/toybox/logger"
)

func Recover(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := recover(); err != nil {
			switch e := err.(type) {
			case error:
				log.Log(c, logger.LevelError, e.Error())
				res.Failed(c, e.Error())
			default:
				log.Log(c, logger.LevelError, "InternalServerError")
				res.Failed(c, "InternalServerError")
			}
			c.Abort()
		}
	}
}
