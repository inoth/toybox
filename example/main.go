package main

import (
	"github.com/inoth/toybox"
	"github.com/inoth/toybox/config"
	"github.com/inoth/toybox/ginsvr"
)

func newApp(conf config.ConfigMate, uc ginsvr.Handler) *toybox.ToyBox {
	t := toybox.New(
		toybox.WithConfig(conf),
		toybox.WithServer(ginsvr.New(
			ginsvr.WithHandlers(uc),
		)),
	)
	return t
}

func main() {
	initApp()
	// app, cleanup, err := initApp()
	// if err != nil {
	// 	panic(err)
	// }
	// defer cleanup()
	// // start and wait for stop signal
	// if err := app.Run(); err != nil && err != context.Canceled {
	// 	panic(err)
	// }
}
