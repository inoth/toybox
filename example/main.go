package main

import (
	"context"
	"os"
)

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
