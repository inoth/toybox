package provider

import (
	"fmt"

	"github.com/inoth/toybox/wsserver"
)

func NewWebsocketServer() *wsserver.WebsocketServer {
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
		),
	)
}
