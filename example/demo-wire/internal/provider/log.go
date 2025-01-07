package provider

import (
	"github.com/inoth/toybox/config"
	"github.com/inoth/toybox/logger"
	"github.com/inoth/toybox/logger/zap"
)

func NewLogger(conf config.ConfigMate) logger.Logger {
	return zap.NewLogger(conf)
}
