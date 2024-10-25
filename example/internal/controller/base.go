package controller

import (
	"example/internal/controller/uqs"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	// NewUserController,
	// NewProxyController,
	// ws.NewMessageController,
	uqs.NewMessageController,
)
