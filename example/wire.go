//go:build wireinject

package main

import (
	"example/internal/controller"
	"example/internal/server"
	"example/internal/service"

	"github.com/google/wire"
	"github.com/inoth/toybox"
	"github.com/inoth/toybox/config"
	"github.com/inoth/toybox/ginsvr"
	"github.com/inoth/toybox/wssvr"
)

func newApp(conf config.ConfigMate, hs *ginsvr.GinHttp2Server, ws *wssvr.WebsocketServer) *toybox.ToyBox {
	t := toybox.New(
		toybox.WithConfig(conf),
		toybox.WithServer(hs, ws),
	)
	return t
}

func initApp(cfg config.CfgBasic) *toybox.ToyBox {
	// panic(wire.Build(config.NewConfig, database.NewDB, service.ProviderSet, controller.ProviderSet, server.ProviderSet, newApp))
	panic(wire.Build(config.NewConfig, service.ProviderSet, controller.ProviderSet, server.ProviderSet, newApp))
}
