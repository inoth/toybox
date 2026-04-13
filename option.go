package toybox

import (
	"os"

	"github.com/inoth/toybox/bootstrap"
	"github.com/inoth/toybox/conf"
	"github.com/inoth/toybox/registry"
	"github.com/inoth/toybox/transport"
)

type Option func(opt *option)

type option struct {
	sigs           []os.Signal
	transports     []transport.Transport
	cfg            conf.ConfigMate
	registrar      registry.Registrar
	serviceName    string
	serviceVersion string
	metadata       map[string]string
	bootstrap      *bootstrap.Config
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

func WithRegistrar(r registry.Registrar) Option {
	return func(opt *option) {
		opt.registrar = r
	}
}

func WithServiceInfo(name, version string) Option {
	return func(opt *option) {
		opt.serviceName = name
		opt.serviceVersion = version
	}
}

func WithMetadata(md map[string]string) Option {
	return func(opt *option) {
		opt.metadata = md
	}
}

// WithBootstrap sets bootstrap config for auto-initializing registry and config source.
// When set, it takes lower priority than explicit WithConfig/WithRegistrar/WithServiceInfo.
func WithBootstrap(cfg *bootstrap.Config) Option {
	return func(opt *option) {
		opt.bootstrap = cfg
	}
}
