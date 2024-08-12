package database

import (
	"fmt"
	"time"

	"github.com/inoth/toybox/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	Name = "mysql"
)

type MysqlComponent struct {
	db *gorm.DB

	Host            string `toml:"host" json:"host"`
	Port            int    `toml:"port" json:"port"`
	User            string `toml:"user" json:"user"`
	Passwd          string `toml:"passwd" json:"passwd"`
	DbName          string `toml:"dbname" json:"dbname"`
	MaxIdleConns    int    `toml:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns    int    `toml:"max_open_conns" json:"max_open_conns"`
	ConnMaxIdletime int    `toml:"conn_max_idletime" json:"conn_max_idletime"`
	ConnMaxLifetime int    `toml:"conn_max_lifetime" json:"conn_max_lifetime"`
}

func NewDB(conf config.ConfigMate) *MysqlComponent {
	mc := MysqlComponent{}
	err := conf.PrimitiveDecode(&mc)
	if err != nil {
		panic(fmt.Errorf("init mysql err: %v", err))
	}
	mc.newDB()
	return &mc
}

func (mc *MysqlComponent) newDB() {
	constr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mc.User, mc.Passwd, mc.Host, mc.Port, mc.DbName)
	client, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       constr,
		DefaultStringSize:         1 << 10,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}))
	if err != nil {
		panic(fmt.Errorf("failed to connect to mysql: %v", err))
	}
	sqlDB, err := client.DB()
	if err != nil {
		panic(fmt.Errorf("failed to get to sql.DB: %v", err))
	}
	sqlDB.SetMaxIdleConns(mc.MaxIdleConns)                                    // 最大空闲连接数
	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(mc.ConnMaxIdletime)) // 连接最大空闲时间
	sqlDB.SetMaxOpenConns(mc.MaxOpenConns)                                    // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(mc.ConnMaxLifetime)) // 连接最大生命周期
	mc.db = client
}

func (mc *MysqlComponent) Name() string {
	return Name
}

func (mc *MysqlComponent) GetDB() any {
	return mc.db
}
