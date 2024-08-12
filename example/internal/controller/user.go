package controller

import (
	"example/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/ginsvr"
)

type UserController struct {
	usvr *service.UserService
}

func (uc *UserController) Prefix() string {
	return "/api"
}

func (uc *UserController) Middlewares() []gin.HandlerFunc {
	return nil
}

func (uc *UserController) Routers() []ginsvr.Router {
	return []ginsvr.Router{
		{Method: "GET", Path: "/user", Handle: []gin.HandlerFunc{uc.UserList}},
	}
}

func NewUserController(usvr *service.UserService) *UserController {
	return &UserController{
		usvr: usvr,
	}
}

func (uc *UserController) UserList(c *gin.Context) {
	c.String(200, "query user list")
}
