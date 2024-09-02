package controller

import (
	"example/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/component/logger"
	"github.com/inoth/toybox/ginsvr"
	"github.com/inoth/toybox/util"
	"github.com/inoth/toybox/util/http/http3"
	"github.com/inoth/toybox/wssvr"
)

type UserController struct {
	log  logger.Logger
	usvr *service.UserService
	hub  *wssvr.WebsocketServer
}

func NewUserController(usvr *service.UserService, hub *wssvr.WebsocketServer) *UserController {
	return &UserController{
		hub:  hub,
		usvr: usvr,
		log:  logger.GetLogger("user_controller"),
	}
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
		{Method: "GET", Path: "/send/:id", Handle: []gin.HandlerFunc{uc.SendMessage}},
		{Method: "GET", Path: "/ws", Handle: []gin.HandlerFunc{uc.Connect}},
		{Method: "GET", Path: "/request", Handle: []gin.HandlerFunc{uc.SendHttp3}},
	}
}

func (uc *UserController) SendHttp3(c *gin.Context) {
	res, err := http3.HttpGet("https://localhost:9060/api/user/1232131231", nil, http3.RequestOption{
		CaCertPath: "cert/ca.pem",
	})
	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.String(200, string(res))
}

func (uc *UserController) UserList(c *gin.Context) {
	uid := c.Param("uid")
	hd := c.GetHeader("SELF_PROXY")
	c.String(200, "key=%s; uid=%s", hd, uid)
}

func (uc *UserController) SendMessage(c *gin.Context) {
	id := c.Param("id")
	msg := c.Query("msg")
	uc.hub.SendMessage([]byte(util.JsonString(map[string]string{"id": id, "body": msg})))
	c.String(200, "ok")
}

func (uc *UserController) Connect(c *gin.Context) {
	clientID, err := wssvr.NewClient(uc.hub, c.Writer, c.Request)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	fmt.Println(clientID)
}
