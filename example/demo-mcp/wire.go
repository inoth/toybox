//go:build wireinject

package main

import (
	"mcp-quickstart/internal/handler"
	"mcp-quickstart/internal/provider"

	"github.com/google/wire"
	"github.com/inoth/toybox"
	"github.com/inoth/toybox/config"
)

func initApp(conf config.ConfigMate) *toybox.ToyBox {
	// panic(wire.Build(data.ProviderSet, biz.ProviderSet, handler.ProviderSet, provider.ProviderSet, newApp))
	panic(wire.Build(handler.ProviderSet, provider.ProviderSet, newApp))
}
