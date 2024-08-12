package main

import (
	"context"

	"github.com/inoth/toybox/config"
)

func main() {
	// initApp()
	app := initApp(config.CfgBasic{
		Remote:   false,
		CfgDir:   "config",
		FileType: "toml",
		Env:      "",
	})

	// start and wait for stop signal
	if err := app.Run(); err != nil && err != context.Canceled {
		panic(err)
	}
}
