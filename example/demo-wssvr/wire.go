//go:build wireinject

package main

import (
	"demo-wssvr/internal/provider"

	"github.com/google/wire"
	"github.com/inoth/toybox"
	"github.com/inoth/toybox/config"
)

func initApp(conf config.ConfigMate) *toybox.ToyBox {
	panic(wire.Build(provider.ProviderSet, newApp))
}
