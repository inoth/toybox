package mysql

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/inoth/toybox/component"

	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	componentName = "mysql"
	mysqlOnce     sync.Once
	Mysql         *MysqlComponents
)

func New() component.Component {
	return &MysqlComponents{}
}

type MysqlComponents struct {
	db *gorm.DB

	Host            string `toml:"host" yaml:"host" json:"host"`
	Port            int    `toml:"port" yaml:"port" json:"port"`
	User            string `toml:"user" yaml:"user" json:"user"`
	Passwd          string `toml:"passwd" yaml:"passwd" json:"passwd"`
	DBName          string `toml:"dbname" yaml:"dbname" json:"dbname"`
	MaxIdleConns    int    `toml:"max_idle_conns" yaml:"max_idle_conns" json:"max_idle_conns"`
	ConnMaxIdleTime int    `toml:"conn_max_idle_time" yaml:"conn_max_idle_time" json:"conn_max_idle_time"`
	MaxOpenConns    int    `toml:"max_open_conns" yaml:"max_open_conns" json:"max_open_conns"`
	ConnMaxLifetime int    `toml:"conn_max_lifetime" yaml:"conn_max_lifetime" json:"conn_max_lifetime"`
}

func (mc *MysqlComponents) Name() string {
	return componentName
}

func (mc *MysqlComponents) String() string {
	buf, _ := json.Marshal(mc)
	return string(buf)
}

func (mc *MysqlComponents) Close() error { return nil }

func (mc *MysqlComponents) Init() (err error) {
	mysqlOnce.Do(func() {
		constr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mc.User, mc.Passwd, mc.Host, mc.Port, mc.DBName)

		mc.db, err = gorm.Open(gmysql.New(gmysql.Config{
			DSN:                       constr,
			DefaultStringSize:         1024, // string 类型字段的默认长度
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}), &gorm.Config{})
		if err != nil {
			err = fmt.Errorf("failed to connect to mysql: %v", err)
			return
		}

		sqlDB, err := mc.db.DB()
		if err != nil {
			return
		}

		// 设置连接池参数
		sqlDB.SetMaxIdleConns(mc.MaxIdleConns)                                    // 最大空闲连接数
		sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(mc.ConnMaxIdleTime)) // 连接最大空闲时间
		sqlDB.SetMaxOpenConns(mc.MaxOpenConns)                                    // 最大打开连接数
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(mc.ConnMaxLifetime)) // 连接最大生命周期

		Mysql = mc
		fmt.Println("mysql component initialization successful")
	})
	return
}

func (mc *MysqlComponents) DB() *gorm.DB {
	return Mysql.db
}
