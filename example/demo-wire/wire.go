//go:build wireinject

package main

import (
	"demo-wire/internal/controller"
	"demo-wire/internal/provider"
	"demo-wire/internal/service"

	"github.com/google/wire"
	"github.com/inoth/toybox"
	"github.com/inoth/toybox/config"
)

func initApp(conf config.ConfigMate) *toybox.ToyBox {
	panic(wire.Build(service.ProviderSet, controller.ProviderSet, provider.ProviderSet, newApp))
}
