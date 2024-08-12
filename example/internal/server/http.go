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
