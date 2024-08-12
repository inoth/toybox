package server

import (
	"example/internal/controller"

	"github.com/inoth/toybox/ginsvr"
)

func NewHttpServer(uc *controller.UserController) *ginsvr.GinHttpServer {
	return ginsvr.New(ginsvr.WithHandlers(uc))
}
