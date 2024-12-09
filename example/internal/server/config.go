package server

import "github.com/inoth/toybox/config"

func NewConfig(dir string) config.ConfigMate {
	return config.NewConfig(
		config.WithConfigDir(dir),
	)
}
