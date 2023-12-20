package logger

type LoggerComponent interface {
	Debug(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
}
