package provider

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/httpserver"
	"github.com/inoth/toybox/wsserver"
)

func NewHttpServer(hub *wsserver.WebsocketServer) *httpserver.GinHttpServer {
	return httpserver.NewHttp(
		httpserver.WithGET("/ws", func(ctx *gin.Context) {
			id, err := wsserver.NewClient(hub, ctx.Writer, ctx.Request)
			if err != nil {
				ctx.String(500, err.Error())
				return
			}
			fmt.Println("new client id:", id)
		}),
	)
}
