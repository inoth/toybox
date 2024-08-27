package server

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	// NewHttp2Server,
	NewHttp3Server,
	NewWebSocketServer,
)
