package mysqlconn

import (
	"fmt"
	"time"

	"github/inoth/ino-toybox/components/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Mysql *gorm.DB

// host: ""
// port: ""
// user: ""
// passwd: ""
// dbname: ""
// max_idle_conns: 10
// max_idle_time: 30
// max_open_conns: 100
// max_life_time: 60
type MysqlComponent struct{}

func (m *MysqlComponent) Init() error {

	host := config.Cfg.GetString("mysql.host")
	port := config.Cfg.GetInt("mysql.port")
	user := config.Cfg.GetString("mysql.user")
	password := config.Cfg.GetString("mysql.passwd")
	database := config.Cfg.GetString("mysql.dbname")

	constr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, database)

	db, err := gorm.Open(mysql.New(mysql.Config{
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
	sqldb.SetMaxIdleConns(config.Cfg.GetInt("mysql.max_idle_conns"))
	sqldb.SetConnMaxIdleTime(time.Duration(config.Cfg.GetInt("mysql.max_idle_time")))
	sqldb.SetMaxOpenConns(config.Cfg.GetInt("mysql.max_open_conns"))
	sqldb.SetConnMaxLifetime(time.Second * time.Duration(config.Cfg.GetInt("mysql.max_life_time")))

	Mysql = db
	return nil
}
