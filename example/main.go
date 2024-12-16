package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/inoth/toybox"
	"github.com/inoth/toybox/config"
	"github.com/inoth/toybox/udpsvr"
)

func newApp(conf config.ConfigMate,
	// hs *ginsvr.GinHttpServer,
	// h2s *ginsvr.GinHttp2Server,
	// h3s *ginsvr.GinHttp3Server,
	// ws *wssvr.WebsocketServer,
	udp *udpsvr.UDPQuicServer,
) *toybox.ToyBox {
	t := toybox.New(
		toybox.WithWatch(),
		toybox.WithConfig(conf),
		toybox.WithServer(
			// hs,
			// h2s,
			// h3s,
			// ws,
			udp,
		),
	)
	return t
}

func main() {
	cfgDir := "config"
	if os.Getenv("CONFIG_ENV") == "dev" {
		cfgDir = "../config"
	}
	cfg := config.NewConfig(
		config.WithConfigDir(cfgDir),
		config.WithConfigInterval(10),
	)
restart:
	app := initApp(cfg)
	// start and wait for stop signal
	if err := app.Run(); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, toybox.ErrRestart) {
		panic(err)
	} else if errors.Is(err, toybox.ErrRestart) {
		log.Println("restart dbproxy, wait 5s...")
		time.Sleep(time.Second * 5)
		goto restart
	}
}
