package httpgin

import (
	"github/inoth/toybox/component/conf"
	"github/inoth/toybox/component/logger"
)

type HttpGinServer struct {
	conf *conf.ConfComponent
	log  logger.LoggerComponent
}

func New(conf *conf.ConfComponent, log logger.LoggerComponent) *HttpGinServer {
	return &HttpGinServer{
		conf: conf,
		log:  log,
	}
}
