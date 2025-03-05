//go:build wireinject

package main

import (
	"demo-wire/internal/biz"
	"demo-wire/internal/data"
	"demo-wire/internal/handler"
	"demo-wire/internal/provider"

	"github.com/google/wire"
	"github.com/inoth/toybox"
	"github.com/inoth/toybox/config"
)

func initApp(conf config.ConfigMate) *toybox.ToyBox {
	panic(wire.Build(data.ProviderSet, biz.ProviderSet, handler.ProviderSet, provider.ProviderSet, newApp))
}
