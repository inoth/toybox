package controller

import (
	"context"
	"example/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/component/logger"
	"github.com/inoth/toybox/ginsvr"
)

type UserController struct {
	log  logger.Logger
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
		{Method: "GET", Path: "/user/:uid", Handle: []gin.HandlerFunc{uc.UserList}},
	}
}

func NewUserController(usvr *service.UserService) *UserController {
	return &UserController{
		usvr: usvr,
		log:  logger.GetLogger("user_controller"),
	}
}

func (uc *UserController) UserList(c *gin.Context) {
	uid := c.Param("uid")
	uc.log.Info(context.Background(), "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	c.JSON(200, uc.usvr.Query(uid))
}
