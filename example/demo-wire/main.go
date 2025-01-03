package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/inoth/toybox"
	"github.com/inoth/toybox/config"
	"github.com/inoth/toybox/config/file"
	"github.com/inoth/toybox/config/toml"
	"github.com/inoth/toybox/ginserver"
	"github.com/inoth/toybox/metric"
	"github.com/inoth/toybox/wsserver"
)

var (
	DefaultDir = "config"
)

func newApp(
	conf config.ConfigMate,
	hs *ginserver.GinHttpServer,
	hs2 *ginserver.GinHttp2Server,
	hs3 *ginserver.GinHttp3Server,
	w *wsserver.WebsocketServer,
	p *metric.Prometheus,
) *toybox.ToyBox {
	t := toybox.New(
		toybox.WithConfig(conf),
		toybox.WithServer(hs),
		toybox.WithServer(hs2),
		toybox.WithServer(hs3),
		toybox.WithServer(w),
		toybox.WithServer(p),
	)
	return t
}

func main() {
	cfg := toml.NewConfiguration(
		config.WithSource(
			file.NewSource(DefaultDir),
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
