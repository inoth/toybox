package server

import (
	"example/internal/controller"

	"github.com/inoth/toybox/ginsvr"
	"github.com/inoth/toybox/validation"
)

func NewHttpServer(uc *controller.UserController) *ginsvr.GinHttpServer {
	return ginsvr.New(
		ginsvr.WithValidator(validation.PhoneValidate),
		ginsvr.WithHandlers(uc),
	)
}

func NewHttp2Server(uc *controller.UserController) *ginsvr.GinHttp2Server {
	return ginsvr.NewHttp2(
		ginsvr.WithValidator(validation.PhoneValidate),
		ginsvr.WithHandlers(uc),
	)
}

func NewHttp3Server(uc *controller.UserController) *ginsvr.GinHttp3Server {
	return ginsvr.NewHttp3(
		ginsvr.WithValidator(validation.PhoneValidate),
		ginsvr.WithHandlers(uc),
	)
}
