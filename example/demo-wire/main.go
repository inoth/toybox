package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/inoth/toybox"
	"github.com/inoth/toybox/config"
	"github.com/inoth/toybox/config/local"
	"github.com/inoth/toybox/config/toml"
	"github.com/inoth/toybox/ginsvr"
	"github.com/inoth/toybox/wssvr"
)

var (
	DefaultDir = "config"
)

func newApp(
	conf config.ConfigMate,
	hs *ginsvr.GinHttpServer,
	hs2 *ginsvr.GinHttp2Server,
	hs3 *ginsvr.GinHttp3Server,
	w *wssvr.WebsocketServer,
) *toybox.ToyBox {
	t := toybox.New(
		toybox.WithConfig(conf),
		toybox.WithServer(hs),
		toybox.WithServer(hs2),
		toybox.WithServer(hs3),
		toybox.WithServer(w),
	)
	return t
}

func main() {
	cfg := toml.NewConfiguration(
		config.WithSource(
			local.NewSource(DefaultDir),
		),
	)

restart:
	app := initApp(cfg)
	if err := app.Run(); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, toybox.ErrRestart) {
		panic(err)
	} else if errors.Is(err, toybox.ErrRestart) {
		log.Println("restart dbproxy, wait 5s...")
		time.Sleep(time.Second * 5)
		goto restart
	}
}
