package controller

import (
	"demo-wire/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/ginsvr"
)

type GreeterController struct {
	svr *service.GreeterService
}

func NewGreeterController(svr *service.GreeterService) *GreeterController {
	return &GreeterController{
		svr: svr,
	}
}

func (gc *GreeterController) Prefix() string {
	return "/api"
}

func (gc *GreeterController) Middlewares() []gin.HandlerFunc {
	return nil
}

func (gc *GreeterController) Routers() []ginsvr.Router {
	return []ginsvr.Router{
		{Method: "GET", Path: "/sayhi/:name", Handle: []gin.HandlerFunc{gc.SayHi}},
	}
}

func (gc *GreeterController) SayHi(c *gin.Context) {
	name := c.Param("name")
	r := gc.svr.SayHi(name)
	c.String(200, r)
}
