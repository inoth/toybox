package main

import (
	"context"
	"errors"

	"github.com/inoth/toybox"
	"github.com/inoth/toybox/config"
	localfile "github.com/inoth/toybox/config/file"
	toml "github.com/inoth/toybox/config/toml"
	"github.com/inoth/toybox/mcpserver"
)

var (
	DefaultDir = "config"
)

func newApp(
	conf config.ConfigMate,
	mcp *mcpserver.MCPServer,
) *toybox.ToyBox {
	t := toybox.New(
		toybox.WithConfig(conf),
		toybox.WithServer(mcp),
	)
	return t
}

func main() {
	cfg := toml.NewConfiguration(
		// config.WithInterval(5),
		config.WithSource(
			localfile.NewLocalSource(DefaultDir),
		),
	)

	app := initApp(cfg)
	if err := app.Run(); err != nil && !errors.Is(err, context.Canceled) {
		panic(err)
	}
}
