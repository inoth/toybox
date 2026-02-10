package toybox

import (
	"os"

	"github.com/inoth/toybox/conf"
	"github.com/inoth/toybox/transport"
)

type Option func(opt *option)

type option struct {
	sigs       []os.Signal
	transports []transport.Transport
	cfg        conf.ConfigMate
}

func WithServer(transport transport.Transport) Option {
	return func(opt *option) {
		opt.transports = append(opt.transports, transport)
	}
}

func WithConfig(cfg conf.ConfigMate) Option {
	return func(opt *option) {
		opt.cfg = cfg
	}
}
