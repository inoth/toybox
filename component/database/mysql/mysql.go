package database

import (
	"fmt"
	"time"

	"github.com/inoth/toybox/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	Name = "mysqls"
)

type Config struct {
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

type MysqlComponent struct {
	DbMap   map[string]*gorm.DB `toml:"-"`
	Configs []Config            `toml:"configs" json:"configs"`
}

func NewGormDatabase(conf config.ConfigMate) *MysqlComponent {
	mc := MysqlComponent{
		DbMap: make(map[string]*gorm.DB),
	}
	err := conf.PrimitiveDecode(&mc)
	if err != nil {
		panic(fmt.Errorf("init mysql err: %v", err))
	}
	mc.initConnect()
	return &mc
}

func (mc *MysqlComponent) Name() string {
	return Name
}

func (mc *MysqlComponent) GetDatabase(dbname string, dst ...any) *gorm.DB {
	if db, ok := mc.DbMap[dbname]; ok {
		if len(dst) > 0 {
			_ = db.AutoMigrate(dst...)
		}
		return db
	}
	panic(fmt.Errorf("not found database %s", dbname))
}

func (mc *MysqlComponent) initConnect() {
	for _, cfg := range mc.Configs {
		constr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User,
			cfg.Passwd,
			cfg.Host,
			cfg.Port,
			cfg.DbName,
		)
		client, err := gorm.Open(mysql.New(mysql.Config{
			DSN:                       constr,
			DefaultStringSize:         1 << 10,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		}))
		if err != nil {
			panic(fmt.Errorf("failed to connect to mysql: %v", err))
		}
		sqlDB, err := client.DB()
		if err != nil {
			panic(fmt.Errorf("failed to get to sql.DB: %v", err))
		}
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)                                    // 最大空闲连接数
		sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(cfg.ConnMaxIdletime)) // 连接最大空闲时间
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)                                    // 最大打开连接数
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(cfg.ConnMaxLifetime)) // 连接最大生命周期

		mc.DbMap[cfg.DbName] = client
	}
}
