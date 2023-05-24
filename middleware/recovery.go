package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/components/logger"
	"github.com/inoth/toybox/res"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.Error(err)
				switch e := err.(type) {
				case error:
					res.Err(c, res.StatusInternalServerError, e)
				default:
					res.Err(c, res.StatusInternalServerError, fmt.Errorf("unknown error %v", err))
				}
			}
		}()
		c.Next()
	}
}
