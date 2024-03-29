package router

import (
	"errors"

	"github.com/inoth/toybox/component/logger"
	"github.com/inoth/toybox/server/ginsvr"
	"github.com/inoth/toybox/server/ginsvr/httpgin"
	"github.com/inoth/toybox/server/ginsvr/validaton"
	"github.com/inoth/toybox/server/metric"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/inoth/toybox/server/ginsvr/res"

	"github.com/gin-gonic/gin"
)

type RequestUser struct {
	Phone string `json:"phone" binding:"emailValidate"`
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
		hgs.GET("", func(c *gin.Context) {
			// if requests_total := metric.GetCounterVec("requests_total"); requests_total != nil {
			// 	requests_total.WithLabelValues(
			// 		"200",
			// 		c.Request.Method,
			// 		c.HandlerName(),
			// 		c.Request.Host,
			// 		c.Request.RequestURI).Inc()
			// }
			metric.CallCounterVec("requests_total", func(cv *prometheus.CounterVec) {
				cv.WithLabelValues("200", c.Request.Method, c.HandlerName(), c.Request.Host, c.Request.RequestURI).Inc()
			})
			res.Ok(c, "ok")
		})
		hgs.POST("", func(ctx *gin.Context) {
			// var user RequestUser
			// if err := ctx.ShouldBindJSON(&user); err != nil {
			// 	res.ErrParamsWithErr(ctx, validaton.ValidatorTranslate(err))
			// 	return
			// }
			user, ok := ginsvr.ParseJsonParam[RequestUser](ctx, func(ctx *gin.Context, err error) {
				res.ErrParamsWithErr(ctx, validaton.ValidatorTranslate(err))
			})
			if !ok {
				return
			}
			res.Ok(ctx, "user", &user)
		})
	}
}
