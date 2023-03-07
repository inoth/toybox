package logger

var Log LoggerComponent

type LoggerComponent interface {
	Info(msg string)
	Infof(msg string, args ...interface{})

	Warn(msg string)
	Warnf(msg string, args ...interface{})

	Err(msg string)
	Errf(msg string, args ...interface{})
}
