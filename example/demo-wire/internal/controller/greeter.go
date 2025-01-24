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

	clientIds []string
}

func NewGreeterController(
	svr *service.GreeterService,
	hub *wsserver.WebsocketServer,
	log logger.Logger,
) *GreeterController {
	return &GreeterController{
		svr:       svr,
		hub:       hub,
		log:       log,
		clientIds: make([]string, 0),
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
		{Method: "POST", Path: "/hi/:name", Handle: []gin.HandlerFunc{gc.SayHi}},
		{Method: "GET", Path: "/ws", Handle: []gin.HandlerFunc{gc.Connect}},
		{Method: "GET", Path: "/clients", Handle: []gin.HandlerFunc{gc.GetClients}},
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

	res.Ok(c, r, gin.H{
		"trace_id": c.Value("trace_id"),
	})
}

func (gc *GreeterController) GetClients(c *gin.Context) {
	res.Ok(c, "", gc.clientIds)
}

func (gc *GreeterController) Connect(c *gin.Context) {
	clientID, err := wsserver.NewClientWithEvent(gc.hub, c.Writer, c.Request,
		func(c *wsserver.Client) {
			gc.clientIds = append(gc.clientIds, c.ID)
		},
		func(c *wsserver.Client) {
			for i, id := range gc.clientIds {
				if c.ID == id {
					gc.clientIds = append(gc.clientIds[:i], gc.clientIds[i+1:]...)
					return
				}
			}
		})
	if err != nil {
		c.String(200, err.Error())
		return
	}
	fmt.Println(clientID)
}
