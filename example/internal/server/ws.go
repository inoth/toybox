package server

import (
	"example/internal/controller/ws"

	"github.com/inoth/toybox/wssvr"
)

func NewWebSocketServer(msg *ws.MessageController) *wssvr.WebsocketServer {
	ws := wssvr.New(wssvr.WithHandler(
		// func(c *wssvr.Context) {
		// 	fmt.Printf("%v\n", string(c.Body()))
		// },
		msg.Handler(),
	))
	return ws
}
