package provider

import (
	"demo-wire/internal/controller/ws"
	"fmt"

	"github.com/inoth/toybox/wssvr"
)

func NewWebsocketServer(
	mc *ws.MessageController,
) *wssvr.WebsocketServer {
	return wssvr.New(
		wssvr.WithHandler(
			func(c *wssvr.Context) {
				fmt.Println("mid 1 start")
				c.Next()
				fmt.Println("mid 1 end")
			},
			func(c *wssvr.Context) {
				fmt.Println("mid 2 start")
				c.Next()
				fmt.Println("mid 2 end")
			},
			func(c *wssvr.Context) {
				fmt.Println("mid 3 start")
				c.Next()
				fmt.Println("mid 3 end")
			},
			mc.Handler(),
		),
	)
}
