package udpsvr

import (
	"context"
	"crypto/tls"
	"time"
)

type Option func(opt *option)

type option struct {
	serverName string

	Addr     string `toml:"addr"`
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`

	ReadBufferSize  int64 `toml:"read_buffer_size"`
	WriteBufferSize int64 `toml:"write_buffer_size"`
	ChannelSize     int64 `toml:"channel_size"`

	WriteWait      time.Duration `toml:"write_wait"`
	PongWait       time.Duration `toml:"pong_wait"`
	PingPeriod     time.Duration `toml:"ping_period"`
	MaxMessageSize int64         `toml:"max_message_size"`
	GZIP           bool          `toml:"gzip"`

	ctx     context.Context
	handles []HandlerFunc
}

func WithContext(ctx context.Context) Option {
	return func(opt *option) {
		opt.ctx = ctx
	}
}

func WithHandler(handles ...HandlerFunc) Option {
	return func(opt *option) {
		opt.handles = append(opt.handles, handles...)
	}
}

func WithTLSFile(certFile, keyFile string) Option {
	return func(opt *option) {
		opt.CertFile = certFile
		opt.KeyFile = keyFile
	}
}

func generateTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-protocol"},
	}, nil
}
