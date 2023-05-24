package logger

var (
	Log ToyBoxLogger
)

type ToyBoxLogger interface {
	Debug(args ...interface{})
	Debugf(msg string, args ...interface{})

	Info(args ...interface{})
	Infof(msg string, args ...interface{})

	Warn(args ...interface{})
	Warnf(msg string, args ...interface{})

	Error(args ...interface{})
	Errorf(msg string, args ...interface{})
}
