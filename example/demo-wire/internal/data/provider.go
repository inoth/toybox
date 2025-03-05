package data

import (
	"demo-wire/internal/data/user"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(user.NewUserRepo)
