package mysql

import (
	"fmt"
	"time"

	"github.com/inoth/ino-toybox/components/config"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Mysql *gorm.DB

// Mysql:
//   Host: ""
//   Port: ""
//   User: ""
//   Passwd: ""
//   DbName: ""
//   MaxIdleConns: 10
//   MaxIdleTime: 30
//   MaxOpenConns: 100
//   MaxLifeTime: 60
type MysqlComponent struct{}

func (m *MysqlComponent) Init() error {

	host := config.Cfg.GetString("Mysql.Host")
	port := config.Cfg.GetInt("Mysql.Port")
	user := config.Cfg.GetString("Mysql.User")
	password := config.Cfg.GetString("Mysql.Passwd")
	database := config.Cfg.GetString("Mysql.DbName")

	constr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, database)

	db, err := gorm.Open(gmysql.New(gmysql.Config{
		DSN:                       constr,
		DefaultStringSize:         1024, // string 类型字段的默认长度
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return err
	}

	sqldb, err := db.DB()
	if err != nil {
		return err
	}
	sqldb.SetMaxIdleConns(config.Cfg.GetInt("Mysql.MaxIdleConns"))
	sqldb.SetConnMaxIdleTime(time.Duration(config.Cfg.GetInt("Mysql.MaxIdleConns")))
	sqldb.SetMaxOpenConns(config.Cfg.GetInt("Mysql.MaxOpenConns"))
	sqldb.SetConnMaxLifetime(time.Second * time.Duration(config.Cfg.GetInt("Mysql.MaxLifeTime")))

	Mysql = db
	return nil
}
