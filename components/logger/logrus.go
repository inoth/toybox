package logger

import (
	"os"

	"github.com/inoth/ino-toybox/components/config"
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Entry

// LogrusLog:
//   LogPath: log/log.log
//   Maxsize: 100
//   MaxAge: 15
//   MaxBackup: 30
//   Compress: true
type LogrusComponent struct {
	hooks         []logrus.Hook
	defaultFields logrus.Fields
	formatter     logrus.Formatter
}

func (lc *LogrusComponent) Init() error {
	log := logrus.New()
	for _, hook := range lc.hooks {
		log.AddHook(hook)
	}
	log.SetReportCaller(true)
	log.SetLevel(logrus.InfoLevel)

	debug := os.Getenv("GORUNEVN")
	if debug == "dev" {
		log.SetFormatter(&logrus.TextFormatter{})
		log.SetOutput(os.Stdout)
		Log = log.WithFields(lc.defaultFields)
		return nil
	}

	if lc.formatter != nil {
		log.SetFormatter(lc.formatter)
	} else {
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	logPath := config.Cfg.GetString("LogrusLog.LogPath")
	logSize := config.Cfg.GetInt("LogrusLog.Maxsize")
	logAge := config.Cfg.GetInt("LogrusLog.MaxAge")
	logBackup := config.Cfg.GetInt("LogrusLog.MaxBackup")
	logCompress := config.Cfg.GetBool("LogrusLog.Compress")

	writer := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    logSize,
		MaxAge:     logAge,
		MaxBackups: logBackup,
		Compress:   logCompress,
	}
	// writerAll := io.MultiWriter(errWriter, warnWriter, infoWriter)
	log.SetOutput(writer)

	Log = log.WithFields(lc.defaultFields)
	return nil
}

func (lc *LogrusComponent) SetHooks(hooks ...logrus.Hook) *LogrusComponent {
	lc.hooks = append(lc.hooks, hooks...)
	return lc
}

func (lc *LogrusComponent) SetSelfFormatter(formatter logrus.Formatter) *LogrusComponent {
	lc.formatter = formatter
	return lc
}

func (lc *LogrusComponent) SetDefaultFields(fields map[string]interface{}) *LogrusComponent {
	lc.defaultFields = fields
	return lc
}
