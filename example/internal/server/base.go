package server

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	// NewHttpServer,
	// NewHttp2Server,
	// NewHttp3Server,
	// NewWebSocketServer,
	NewUDPQuicServer,
)
