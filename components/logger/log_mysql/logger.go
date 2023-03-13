package logmysql

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/inoth/toybox/components/config"
	"github.com/inoth/toybox/components/database/mysqlconn"
	"github.com/inoth/toybox/components/logger"
)

// CREATE TABLE `t_logger` (
//   `id` int(11) NOT NULL AUTO_INCREMENT,
//   `app_name` varchar(100) NOT NULL COMMENT '项目名称',
//   `level` varchar(100) NOT NULL COMMENT '日志等级',
//   `msg` varchar(1024) DEFAULT NULL,
//   `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
//   PRIMARY KEY (`id`)
// ) ENGINE=InnoDB AUTO_INCREMENT=3001 DEFAULT CHARSET=latin1 COMMENT='日志记录表';
type MysqlLogger struct {
	Level string `gorm:"level"`
	Msg   string `gorm:"msg"`
}

func (MysqlLogger) TableName() string {
	return "t_logger"
}

// app_name: datantv1
// max_channel_pool: 1000
type LogMysqlComponent struct {
	appName string
	msg     chan MysqlLogger
	db      *sql.DB
	env     string
}

// 管道、数据库连接、队列
func (lmc *LogMysqlComponent) Init() (err error) {
	lmc.env = os.Getenv("GORUNEVN")
	lmc.appName = config.Cfg.GetString("logger.app_name")

	if lmc.env != "dev" {
		lmc.db, err = mysqlconn.DB.DB()
		if err != nil {
			return err
		}
	}

	// 允许写入数据库协程
	lmc.RunLogger()

	logger.Log = lmc
	return
}

// 设定消息写入池
func (lmc *LogMysqlComponent) RunLogger() {
	// 初始化写入缓存池
	lmc.msg = make(chan MysqlLogger, config.Cfg.GetInt("logger.max_channel_pool"))
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		if lmc.env != "dev" {
			statement, err := lmc.db.Prepare("INSERT INTO t_logger(`app_name`,`level`,`msg`) VALUES(?,?,?)")
			if err != nil {
				panic(err)
			}
			for msg := range lmc.msg {
				_, err = statement.Exec(lmc.appName, msg.Level, msg.Msg)
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			for msg := range lmc.msg {
				fmt.Println(msg.Msg)
			}
		}
	}()
}

func (lmc *LogMysqlComponent) Info(msg string) {
	lmc.msg <- MysqlLogger{
		Level: "info",
		Msg:   msg,
	}
}

func (lmc *LogMysqlComponent) Infof(msg string, args ...interface{}) {
	lmc.msg <- MysqlLogger{
		Level: "info",
		Msg:   fmt.Sprintf(msg, args...),
	}
}

func (lmc *LogMysqlComponent) Warn(msg string) {
	lmc.msg <- MysqlLogger{
		Level: "warn",
		Msg:   msg,
	}
}
func (lmc *LogMysqlComponent) Warnf(msg string, args ...interface{}) {
	lmc.msg <- MysqlLogger{
		Level: "warn",
		Msg:   fmt.Sprintf(msg, args...),
	}
}

func (lmc *LogMysqlComponent) Err(msg string) {
	lmc.msg <- MysqlLogger{
		Level: "error",
		Msg:   msg,
	}
}

func (lmc *LogMysqlComponent) Errf(msg string, args ...interface{}) {
	lmc.msg <- MysqlLogger{
		Level: "error",
		Msg:   fmt.Sprintf(msg, args...),
	}
}
