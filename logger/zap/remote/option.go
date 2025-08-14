package logger

import "io"

type Option func(*ZapRemoteLogger)

func WithWriter(w io.Writer) Option {
	return func(z *ZapRemoteLogger) {
		z.writer = w
	}
}
