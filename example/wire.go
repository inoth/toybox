//go:build wireinject

package main

import (
	"example/internal/controller"
	"example/internal/server"
	"example/internal/service"

	"github.com/google/wire"
	"github.com/inoth/toybox"
	"github.com/inoth/toybox/component/database"
	"github.com/inoth/toybox/config"
	"github.com/inoth/toybox/ginsvr"
)

func newApp(conf config.ConfigMate, hs *ginsvr.GinHttpServer) *toybox.ToyBox {
	t := toybox.New(
		toybox.WithConfig(conf),
		toybox.WithServer(hs),
	)
	return t
}

func initApp(cfg config.CfgBasic) *toybox.ToyBox {
	panic(wire.Build(config.NewConfig, database.NewDB, service.ProviderSet, controller.ProviderSet, server.ProviderSet, newApp))
	// panic(wire.Build(config.NewConfig, service.ProviderSet, controller.ProviderSet, server.ProviderSet, newApp))
}
