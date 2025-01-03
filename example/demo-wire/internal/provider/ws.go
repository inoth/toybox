package provider

import (
	"demo-wire/internal/controller/ws"
	"fmt"

	"github.com/inoth/toybox/wsserver"
)

func NewWebsocketServer(
	mc *ws.MessageController,
) *wsserver.WebsocketServer {
	return wsserver.New(
		wsserver.WithHandler(
			func(c *wsserver.Context) {
				fmt.Println("mid 1 start")
				c.Next()
				fmt.Println("mid 1 end")
			},
			func(c *wsserver.Context) {
				fmt.Println("mid 2 start")
				c.Next()
				fmt.Println("mid 2 end")
			},
			func(c *wsserver.Context) {
				fmt.Println("mid 3 start")
				c.Next()
				fmt.Println("mid 3 end")
			},
			mc.Handler(),
		),
	)
}
