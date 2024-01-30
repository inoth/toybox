//go:build wireinject

package main

import (
	"demo/handler"
	"demo/service"

	"github.com/inoth/toybox"
	"github.com/inoth/toybox/component/logger"
	"github.com/inoth/toybox/server/ginsvr/httpgin"
	"github.com/inoth/toybox/server/websocket"

	"github.com/google/wire"
)

func newApp(conf toybox.ConfigMate, foo *handler.FooHandler) *toybox.ToyBox {
	return toybox.New(
		logger.WithLogger(),
		toybox.WithConfigMate(conf),
		websocket.NewWebsocketServer(websocket.WithMessageHandler(func(c *websocket.Context) {
			input := c.GetMessage()
			c.SendMessage(websocket.NewOutputMessage([]byte("message handler with "+input.Id), input.Id))
		})),
		httpgin.NewHttpGin(httpgin.WithHandlers(foo)),
	)
}

func initApp() *toybox.ToyBox {
	panic(wire.Build(toybox.NewConfig, service.NewFooService, handler.NewFooHandler, newApp))
}
