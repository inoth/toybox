package provider

import (
	"github.com/google/wire"
	database "github.com/inoth/toybox/component/database/sqlite"
)

var ProviderSet = wire.NewSet(
	NewHttpServer,
	// NewHttp2Server,
	// NewHttp3Server,
	// NewWebsocketServer,
	NewMetric,
	NewProperty,
	NewLogger,
	database.NewGormDatabase,
)
