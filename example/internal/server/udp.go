package server

import (
	"example/internal/controller/uqs"

	"github.com/inoth/toybox/udpsvr"
)

func NewUDPQuicServer(col *uqs.MessageController) *udpsvr.UDPQuicServer {
	ws := udpsvr.New(udpsvr.WithHandler(
		// func(c *udpsvr.Context) {
		// 	fmt.Printf("%v\n", string(c.Body()))
		// },
		col.Handler(),
	))
	return ws
}
