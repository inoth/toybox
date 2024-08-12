// go:build wireinject

package main

import (
	"example/internal/controller"

	"github.com/google/wire"
	"github.com/inoth/toybox"
	"github.com/inoth/toybox/config"
)

// func newApp(conf config.ConfigMate, uc ginsvr.Handler) *toybox.ToyBox {
// 	t := toybox.New(
// 		toybox.WithConfig(conf),
// 		toybox.WithServer(ginsvr.New(
// 			ginsvr.WithHandlers(uc),
// 		)),
// 	)
// 	return t
// }

func initApp() *toybox.ToyBox {
	cfg := config.NewConfig(config.CfgBasic{
		Remote:   false,
		CfgDir:   "config",
		FileType: "toml",
		Env:      "",
	})
	panic(wire.Build(cfg, controller.ControllerSet, newApp))
	return nil
}
