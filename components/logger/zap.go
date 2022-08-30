package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/inoth/ino-toybox/components/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Zap *zap.Logger

type ZapLogger struct {
	hooks []func(zapcore.Entry) error
}

func (zpl *ZapLogger) SetHooks(hooks ...func(zapcore.Entry) error) *ZapLogger {
	zpl.hooks = hooks
	return zpl
}

func (zpl *ZapLogger) Init() error {
	Zap = newLogger(zpl.hooks...)
	return nil
}

// NewLogger 创建自定义zap logger对象
func newLogger(hooks ...func(zapcore.Entry) error) *zap.Logger {
	// 若为debug模式，创建debug日志级别的logger对象，直接输出到屏幕，不写入文件
	debug := os.Getenv("GORUNEVN")
	if debug == "dev" {
		logger, err := zap.NewDevelopment(zap.Hooks(hooks...))
		if err != nil {
			fmt.Printf("创建zap日志包失败，详情：%s\n", err.Error())
		}
		return logger
	}

	// 解析配置文件，获取日志相关配置
	errLog := config.Cfg.GetString("ZapLog.ErrLog")
	warnLog := config.Cfg.GetString("ZapLog.WarnLog")
	infoLog := config.Cfg.GetString("ZapLog.InfoLog")
	logSize := config.Cfg.GetInt("ZapLog.Maxsize")
	logAge := config.Cfg.GetInt("ZapLog.MaxAge")
	logBackup := config.Cfg.GetInt("ZapLog.MaxBackup")
	logCompress := config.Cfg.GetBool("ZapLog.Compress")
	jsonFormat := config.Cfg.GetBool("ZapLog.Json")

	// 设置日志内容格式，以及日志输出格式。默认为人类可读格式；若配置了json，则输出为json格式
	encoderConf := genEncoderConf()
	encoder := zapcore.NewConsoleEncoder(encoderConf)
	if jsonFormat {
		encoder = zapcore.NewJSONEncoder(encoderConf)
	}

	// 错误日志
	errWriter := &lumberjack.Logger{
		Filename:   errLog,
		MaxSize:    logSize,
		MaxAge:     logAge,
		MaxBackups: logBackup,
		Compress:   logCompress,
	}

	// 警告日志
	warnWriter := &lumberjack.Logger{
		Filename:   warnLog,
		MaxSize:    logSize,
		MaxAge:     logAge,
		MaxBackups: logBackup,
		Compress:   logCompress,
	}
	// 普通日志
	infoWriter := &lumberjack.Logger{
		Filename:   infoLog,
		MaxSize:    logSize,
		MaxAge:     logAge,
		MaxBackups: logBackup,
		Compress:   logCompress,
	}

	// 日志级别配置，不能直接写zap.InfoLevel等，否则在写error级别的log时，info、warn也会写一份
	errLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zap.ErrorLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv > zap.InfoLevel && lv <= zap.WarnLevel
	})
	infoLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv > zap.DebugLevel && lv <= zap.InfoLevel
	})

	// 启用多个输出流，不同级别的日志写到不同的日志文件中
	// 由于启用了多个输出流，所以配置文件中不必设置log_level，没有意义
	writers := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.AddSync(errWriter), errLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
	}

	// 创建zap logger对象，同时添加两个option：日志打印行号、自定义hook
	return zap.New(zapcore.NewTee(writers...), zap.AddCaller(), zap.AddCallerSkip(1), zap.Hooks(hooks...))
}

// genEncoderConf 生成EncoderConfig，用于配置日志格式
func genEncoderConf() zapcore.EncoderConfig {
	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.EncodeTime = zapTimeEncoder               // 日志规范要求时间格式到毫秒
	encoderConf.TimeKey = "created_at"                    // 时间戳的key使用timestamp，根据model的定义进行设置
	encoderConf.MessageKey = "message"                    // 消息的key使用message
	encoderConf.EncodeLevel = zapcore.CapitalLevelEncoder // 日志规范要求日志级别为大写格式
	return encoderConf
}

// zapTimeEncoder 用于日志时间格式化，到毫秒级
func zapTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05:000"))
}
