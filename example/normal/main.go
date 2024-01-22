package main

import (
	"demo1/router"

	"github.com/inoth/toybox"
	"github.com/inoth/toybox/component/logger"
	"github.com/inoth/toybox/server/ginsvr"
	"github.com/inoth/toybox/server/ginsvr/httpgin"
)

func main() {
	tb := toybox.New(
		toybox.WithLoadConf(),
		logger.WithLogger(),
		httpgin.NewHttpGin(func(hgs *httpgin.HttpGinServer) {
			hgs.Use(ginsvr.Recovery())
		}, router.WithUserRouter()),
	)
	if err := tb.Run(); err != nil {
		panic(err)
	}
}
