package provider

import (
	"demo-wire/internal/controller"

	"github.com/inoth/toybox/ginserver"
	"github.com/inoth/toybox/ginserver/middleware"
)

func NewHttpServer(
	gc *controller.GreeterController,
) *ginserver.GinHttpServer {
	return ginserver.NewHttp(
		ginserver.WithMiddleware(middleware.SetTraceId()),
		ginserver.WithHandlers(gc),
	)
}

func NewHttp2Server(
	gc *controller.GreeterController,
) *ginserver.GinHttp2Server {
	return ginserver.NewHttp2(
		ginserver.WithMiddleware(middleware.SetTraceId()),
		ginserver.WithHandlers(gc),
	)
}

func NewHttp3Server(
	gc *controller.GreeterController,
) *ginserver.GinHttp3Server {
	return ginserver.NewHttp3(
		ginserver.WithMiddleware(middleware.SetTraceId()),
		ginserver.WithHandlers(gc),
	)
}
