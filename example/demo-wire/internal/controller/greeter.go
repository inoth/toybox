package controller

import (
	"demo-wire/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/ginsvr"
	"github.com/inoth/toybox/wssvr"
)

type GreeterController struct {
	svr *service.GreeterService
	hub *wssvr.WebsocketServer
}

func NewGreeterController(svr *service.GreeterService, hub *wssvr.WebsocketServer) *GreeterController {
	return &GreeterController{
		svr: svr,
		hub: hub,
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
		{Method: "GET", Path: "/ws", Handle: []gin.HandlerFunc{gc.Connect}},
	}
}

func (gc *GreeterController) SayHi(c *gin.Context) {
	name := c.Param("name")
	r := gc.svr.SayHi(name)
	c.String(200, r)
}

func (uc *GreeterController) Connect(c *gin.Context) {
	clientID, err := wssvr.NewClient(uc.hub, c.Writer, c.Request)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	fmt.Println(clientID)
}
