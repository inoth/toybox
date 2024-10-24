package server

import (
	"fmt"

	"github.com/inoth/toybox/udpsvr"
)

func NewUDPQuicServer() *udpsvr.UDPQuicServer {
	ws := udpsvr.New(udpsvr.WithHandler(
		func(c *udpsvr.Context) {
			fmt.Printf("%v\n", string(c.Body()))
		},
	))
	return ws
}
