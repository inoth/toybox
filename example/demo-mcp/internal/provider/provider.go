package provider

import (
	"github.com/google/wire"
	database "github.com/inoth/toybox/component/database/sqlite"
)

var ProviderSet = wire.NewSet(
	NewLogger,
	NewMCPServer,
	database.NewGormDatabase,
)
