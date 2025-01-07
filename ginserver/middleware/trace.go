package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/util"
)

func SetTraceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := util.UUID(32)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "trace_id", traceId))
		c.Request.Header.Add("trace_id", traceId)
		c.Next()
	}
}
