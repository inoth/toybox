package router

import (
	"errors"

	"github.com/inoth/toybox/component/logger"
	"github.com/inoth/toybox/server/ginsvr/httpgin"
	"github.com/inoth/toybox/server/ginsvr/validators"

	"github.com/inoth/toybox/server/ginsvr/res"

	"github.com/gin-gonic/gin"
)

type RequestUser struct {
	Phone string `json:"phone" binding:"phoneValidate"`
}

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
		hgs.POST("", func(ctx *gin.Context) {
			var user RequestUser
			if err := ctx.BindJSON(&user); err != nil {
				res.ErrParamsWithErr(ctx, validators.ValidatorTranslate(err))
				return
			}
			res.Ok(ctx, "user", &user)
		})
	}
}
