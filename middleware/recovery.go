package middleware

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/inoth/ino-toybox/res"
	"github.com/inoth/ino-toybox/utils"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if utils.GetRunEnv() == utils.EnvDev {
					fmt.Println(string(debug.Stack()))
				}
				switch e := err.(type) {
				case error:
					res.ResultErr(c, res.InternalErrorCode, e)
				default:
					res.ResultErr(c, res.InternalErrorCode, errors.New("internal server error"))
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}
