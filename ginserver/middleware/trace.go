package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/util"
)

func SetTraceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Header.Add("trace_id", util.UUID(32))
		c.Next()
	}
}
