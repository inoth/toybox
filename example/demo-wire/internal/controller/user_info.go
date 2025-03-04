package controller

import (
	"demo-wire/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/ginserver"
	"github.com/inoth/toybox/ginserver/res"
)

type UserInfoController struct {
	svr *service.UserInfoService
}

func NewUserInfoService(svr *service.UserInfoService) *UserInfoController {
	return &UserInfoController{svr: svr}
}

func (gc *UserInfoController) Prefix() string {
	return "/api/user"
}

func (gc *UserInfoController) Middlewares() []gin.HandlerFunc {
	return nil
}

func (gc *UserInfoController) Routers() []ginserver.Router {
	return []ginserver.Router{
		{Method: "GET", Path: "get/:name", Handle: []gin.HandlerFunc{gc.GetUserInfoByName}},
		{Method: "GET", Path: "set/:name", Handle: []gin.HandlerFunc{gc.CreateUserInfo}},
	}
}

func (gc *UserInfoController) GetUserInfoByName(c *gin.Context) {
	name := c.Param("name")
	userInfo, err := gc.svr.GetUserInfoByName(name)
	if err != nil {
		res.ErrParams(c, name)
		return
	}
	res.Ok(c, "ok", userInfo)
}

func (gc *UserInfoController) CreateUserInfo(c *gin.Context) {
	name := c.Param("name")
	userInfo, err := gc.svr.CreateUserInfo(name)
	if err != nil {
		res.ErrParams(c, name)
		return
	}
	res.Ok(c, "ok", userInfo)
}
