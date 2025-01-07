package controller

import (
	"demo-wire/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/ginserver"
	"github.com/inoth/toybox/ginserver/res"
	"github.com/inoth/toybox/logger"
	"github.com/inoth/toybox/wsserver"
)

type GreeterController struct {
	svr *service.GreeterService
	hub *wsserver.WebsocketServer
	log logger.Logger
}

func NewGreeterController(
	svr *service.GreeterService,
	hub *wsserver.WebsocketServer,
	log logger.Logger,
) *GreeterController {
	return &GreeterController{
		svr: svr,
		hub: hub,
		log: log,
	}
}

func (gc *GreeterController) Prefix() string {
	return "/api"
}

func (gc *GreeterController) Middlewares() []gin.HandlerFunc {
	return nil
}

func (gc *GreeterController) Routers() []ginserver.Router {
	return []ginserver.Router{
		{Method: "GET", Path: "/sayhi/:name", Handle: []gin.HandlerFunc{gc.SayHi}},
		{Method: "GET", Path: "/ws", Handle: []gin.HandlerFunc{gc.Connect}},
	}
}

func (gc *GreeterController) SayHi(c *gin.Context) {
	name := c.Param("name")
	r := gc.svr.SayHi(name)

	// logger.Log(c, logger.LevelInfo, "this is info logger")
	gc.log.Log(c, logger.LevelDebug, "this is debug logger")
	gc.log.Log(c, logger.LevelInfo, "this is info logger")
	gc.log.Log(c, logger.LevelWarn, "this is warn logger")
	gc.log.Log(c, logger.LevelError, "this is error logger")

	res.Ok(c, "", gin.H{
		"trace_id": c.Value("trace_id"),
		"msg":      r,
	})
}

func (uc *GreeterController) Connect(c *gin.Context) {
	clientID, err := wsserver.NewClient(uc.hub, c.Writer, c.Request)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	fmt.Println(clientID)
}
