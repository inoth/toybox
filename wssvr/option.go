package wssvr

import (
	"context"
	"time"
)

type Option func(opt *option)

type option struct {
	serverName string

	ReadBufferSize  int64 `toml:"read_buffer_size"`
	WriteBufferSize int64 `toml:"write_buffer_size"`
	ChannelSize     int64 `toml:"channel_size"`

	WriteWait      time.Duration `toml:"write_wait"`
	PongWait       time.Duration `toml:"pong_wait"`
	PingPeriod     time.Duration `toml:"ping_period"`
	MaxMessageSize int64         `toml:"max_message_size"`

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
