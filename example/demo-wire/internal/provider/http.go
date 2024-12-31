package provider

import (
	"demo-wire/internal/controller"

	"github.com/inoth/toybox/ginsvr"
)

func NewHttpServer(
	gc *controller.GreeterController,
) *ginsvr.GinHttpServer {
	return ginsvr.NewHttp(
		ginsvr.WithHandlers(gc),
	)
}

func NewHttp2Server(
	gc *controller.GreeterController,
) *ginsvr.GinHttp2Server {
	return ginsvr.NewHttp2(
		ginsvr.WithHandlers(gc),
	)
}

func NewHttp3Server(
	gc *controller.GreeterController,
) *ginsvr.GinHttp3Server {
	return ginsvr.NewHttp3(
		ginsvr.WithHandlers(gc),
	)
}
