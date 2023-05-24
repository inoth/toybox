package zap

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/inoth/toybox/component"
	"github.com/inoth/toybox/components/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	componentName = "logger"
	logOnce       sync.Once
)

type ZapLogComponent struct {
	zap *zap.Logger

	ErrLog    string `json:"err_log" yaml:"err_log" toml:"err_log"`
	WarnLog   string `json:"warn_log" yaml:"warn_log" toml:"warn_log"`
	InfoLog   string `json:"info_log" yaml:"info_log" toml:"info_log"`
	DebugLog  string `json:"debug_log" yaml:"debug_log" toml:"debug_log"`
	MaxSize   int    `json:"max_size" yaml:"max_size" toml:"max_size"`
	MaxAge    int    `json:"max_age" yaml:"max_age" toml:"max_age"`
	MaxBackup int    `json:"max_backup" yaml:"max_backup" toml:"max_backup"`
	Compress  bool   `json:"compress" yaml:"compress" toml:"compress"`
	Json      bool   `json:"json" yaml:"json" toml:"json"`
}

func New() component.Component {
	return &ZapLogComponent{}
}

func (zc *ZapLogComponent) Name() string   { return componentName }
func (zc *ZapLogComponent) String() string { return "" }

func (zc *ZapLogComponent) Init() (err error) {
	logOnce.Do(func() {
		zc.zap, err = zc.newLogger()
		if err != nil {
			return
		}
		logger.Log = zc
	})
	return
}

func (zc *ZapLogComponent) newLogger() (*zap.Logger, error) {
	debug := os.Getenv("GORUNEVN")
	if debug == "dev" || debug == "debug" {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize the logging component: %v", err)
		}
		return logger, nil
	}

	encoderConf := genEncoderConf()
	encoder := zapcore.NewConsoleEncoder(encoderConf)
	if zc.Json {
		encoder = zapcore.NewJSONEncoder(encoderConf)
	}

	errWriter := &lumberjack.Logger{
		Filename:   zc.ErrLog,
		MaxSize:    zc.MaxSize,
		MaxAge:     zc.MaxAge,
		MaxBackups: zc.MaxBackup,
		Compress:   zc.Compress,
	}

	warnWriter := &lumberjack.Logger{
		Filename:   zc.WarnLog,
		MaxSize:    zc.MaxSize,
		MaxAge:     zc.MaxAge,
		MaxBackups: zc.MaxBackup,
		Compress:   zc.Compress,
	}

	infoWriter := &lumberjack.Logger{
		Filename:   zc.InfoLog,
		MaxSize:    zc.MaxSize,
		MaxAge:     zc.MaxAge,
		MaxBackups: zc.MaxBackup,
		Compress:   zc.Compress,
	}

	errLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zap.ErrorLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv > zap.InfoLevel && lv <= zap.WarnLevel
	})
	infoLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv > zap.DebugLevel && lv <= zap.InfoLevel
	})

	writers := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.AddSync(errWriter), errLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
	}

	return zap.New(zapcore.NewTee(writers...), zap.AddCaller(), zap.AddCallerSkip(1)), nil
}

func genEncoderConf() zapcore.EncoderConfig {
	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.EncodeTime = zapTimeEncoder
	encoderConf.TimeKey = "created_at"
	encoderConf.MessageKey = "message"
	encoderConf.EncodeLevel = zapcore.CapitalLevelEncoder
	return encoderConf
}

func zapTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05:000"))
}

func (zc *ZapLogComponent) Debug(args ...interface{}) {
	suger := zc.zap.Sugar()
	suger.Debug(args...)
}
func (zc *ZapLogComponent) Debugf(msg string, args ...interface{}) {
	suger := zc.zap.Sugar()
	suger.Debugf(msg, args...)
}

func (zc *ZapLogComponent) Info(args ...interface{}) {
	suger := zc.zap.Sugar()
	suger.Info(args...)
}
func (zc *ZapLogComponent) Infof(msg string, args ...interface{}) {
	suger := zc.zap.Sugar()
	suger.Infof(msg, args...)
}

func (zc *ZapLogComponent) Warn(args ...interface{}) {
	suger := zc.zap.Sugar()
	suger.Warn(args...)
}
func (zc *ZapLogComponent) Warnf(msg string, args ...interface{}) {
	suger := zc.zap.Sugar()
	suger.Warnf(msg, args...)
}

func (zc *ZapLogComponent) Error(args ...interface{}) {
	suger := zc.zap.Sugar()
	suger.Error(args...)
}
func (zc *ZapLogComponent) Errorf(msg string, args ...interface{}) {
	suger := zc.zap.Sugar()
	suger.Errorf(msg, args...)
}
