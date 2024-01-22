package router

import (
	"errors"

	"github.com/inoth/toybox/component/logger"
	"github.com/inoth/toybox/server/ginsvr/httpgin"

	"github.com/inoth/toybox/server/ginsvr/res"

	"github.com/gin-gonic/gin"
)

func WithUserRouter() httpgin.Option {
	return func(hgs *httpgin.HttpGinServer) {
		user := hgs.Group("user")
		{
			user.GET("", func(c *gin.Context) {
				panic(errors.New("self error"))
				log := logger.New("user", "")
				log.Debug("user_info")
				log.Info("user_info")
				log.Warn("user_info")
				log.Error("user_info")
				res.Ok(c, "user_info")
			})
		}
	}
}
