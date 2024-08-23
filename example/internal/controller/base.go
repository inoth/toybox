package controller

import (
	"example/internal/controller/ws"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserController,
	ws.NewMessageController,
)
