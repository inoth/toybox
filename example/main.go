package main

import (
	"github/inoth/toybox"
	"github/inoth/toybox/server/httpgin"

	"github.com/gin-gonic/gin"
)

func main() {
	tb := toybox.New(httpgin.NewHttpGin(func(hgs *httpgin.HttpGinServer) {
		hgs.GET("test", func(c *gin.Context) {
			c.String(200, "ok")
		})
	}))
	tb.Run()
}
