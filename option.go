package toybox

import (
	"os"

	"github.com/inoth/toybox/conf"
	"github.com/inoth/toybox/transport"
)

type Option func(opt *option)

type option struct {
	sigs []os.Signal
	svcs []transport.Transport
	cfg  conf.ConfigMate
}

func WithServer(svcs ...transport.Transport) Option {
	return func(opt *option) {
		opt.svcs = append(opt.svcs, svcs...)
	}
}

func WithConfig(cfg conf.ConfigMate) Option {
	return func(opt *option) {
		opt.cfg = cfg
	}
}
