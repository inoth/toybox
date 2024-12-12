package main

import (
	"context"
	"os"

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
	app := initApp(cfgDir)

	// start and wait for stop signal
	if err := app.Run(); err != nil && err != context.Canceled {
		panic(err)
	}
}
