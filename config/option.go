package config

type Options func(opt *Option)

type Option struct {
	Interval int
	Source   Source
}

func WithInterval(interval int) Options {
	return func(opt *Option) {
		opt.Interval = interval
	}
}

func WithSource(source Source) Options {
	return func(opt *Option) {
		opt.Source = source
	}
}
