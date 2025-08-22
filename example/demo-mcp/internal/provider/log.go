package provider

import (
	"github.com/inoth/toybox/config"
	"github.com/inoth/toybox/logger"
	zaplog "github.com/inoth/toybox/logger/zap"
)

func NewLogger(conf config.ConfigMate) logger.Logger {
	return zaplog.NewZapLogger(conf)
}
