package logger

import (
	"go.uber.org/zap"
)

type LoggerConfig struct {
	TraceId    string
	ServerName string
}

type Logger struct {
	LoggerConfig
}

func New(name, traceId string) *Logger {
	return GetLogger(LoggerConfig{ServerName: name, TraceId: traceId})
}

func GetLogger(cfg ...LoggerConfig) *Logger {
	logger := &Logger{}
	if len(cfg) > 0 {
		logger.LoggerConfig = cfg[0]
	}
	return logger
}

func (l *Logger) Debug(msg string) {
	if log == nil {
		return
	}
	log.Debug(msg,
		zap.String("server_name", l.ServerName),
		zap.String("trace_id", l.TraceId),
	)
}

func (l *Logger) Info(msg string) {
	if log == nil {
		return
	}
	log.Info(msg,
		zap.String("server_name", l.ServerName),
		zap.String("trace_id", l.TraceId),
	)
}

func (l *Logger) Warn(msg string) {
	if log == nil {
		return
	}
	log.Warn(msg,
		zap.String("server_name", l.ServerName),
		zap.String("trace_id", l.TraceId),
	)
}

func (l *Logger) Error(msg string) {
	if log == nil {
		return
	}
	log.Error(msg,
		zap.String("server_name", l.ServerName),
		zap.String("trace_id", l.TraceId),
	)
}

func Debug(msg string) {
	if log == nil {
		return
	}
	log.Debug(msg)
}

func Info(msg string) {
	if log == nil {
		return
	}
	log.Info(msg)
}

func Warn(msg string) {
	if log == nil {
		return
	}
	log.Warn(msg)
}

func Error(msg string) {
	if log == nil {
		return
	}
	log.Error(msg)
}
