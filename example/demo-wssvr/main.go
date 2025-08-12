package main

import (
	"context"
	"errors"

	"github.com/inoth/toybox"
	"github.com/inoth/toybox/config"
	localfile "github.com/inoth/toybox/config/file"
	toml "github.com/inoth/toybox/config/toml"
	"github.com/inoth/toybox/httpserver"
	"github.com/inoth/toybox/wsserver"
)

var (
	DefaultDir = "config"
)

func newApp(
	conf config.ConfigMate,
	hs *httpserver.GinHttpServer,
	w *wsserver.WebsocketServer,
) *toybox.ToyBox {
	t := toybox.New(
		toybox.WithConfig(conf),
		toybox.WithServer(hs, w),
	)
	return t
}

func main() {
	cfg := toml.NewConfiguration(
		config.WithSource(
			localfile.NewLocalSource(DefaultDir),
		),
	)

	app := initApp(cfg)
	if err := app.Run(); err != nil && !errors.Is(err, context.Canceled) {
		panic(err)
	}
}
