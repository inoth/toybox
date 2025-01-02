package controller

import (
	"demo-wire/internal/controller/ws"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewGreeterController,
	ws.NewMessageController,
)
