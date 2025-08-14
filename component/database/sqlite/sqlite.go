package database

import (
	"fmt"
	"time"

	"github.com/inoth/toybox/config"
	// "gorm.io/driver/sqlite" // 需要开启CGO

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const (
	name = "sqlites"
)

type Config struct {
	DBPath          string `toml:"db_path" json:"db_path"`
	DBName          string `toml:"dbname" json:"dbname"`
	MaxIdleConns    int    `toml:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns    int    `toml:"max_open_conns" json:"max_open_conns"`
	ConnMaxIdletime int    `toml:"conn_max_idletime" json:"conn_max_idletime"`
	ConnMaxLifetime int    `toml:"conn_max_lifetime" json:"conn_max_lifetime"`
}

type SqliteComponent struct {
	DBMap   map[string]*gorm.DB `toml:"-"`
	Configs []Config            `toml:"configs" json:"configs"`
}

func NewGormDatabase(conf config.ConfigMate) *SqliteComponent {
	sc := SqliteComponent{
		DBMap: make(map[string]*gorm.DB),
	}
	err := conf.PrimitiveDecode(&sc)
	if err != nil {
		panic(fmt.Errorf("init mysql err: %v", err))
	}
	sc.initConnect()
	return &sc
}

func (sc *SqliteComponent) Name() string {
	return name
}

func (sc *SqliteComponent) GetDatabase(dbName string, dst ...any) *gorm.DB {
	if db, ok := sc.DBMap[dbName]; ok {
		if len(dst) > 0 {
			_ = db.AutoMigrate(dst...)
		}
		return db
	}
	panic(fmt.Errorf("not found database %s", dbName))
}

func (sc *SqliteComponent) initConnect() {
	for _, cfg := range sc.Configs {
		client, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("failed to connect to sqlite: %v", err))
		}
		sqlDB, err := client.DB()
		if err != nil {
			panic(fmt.Errorf("failed to get to sql.DB: %v", err))
		}
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)                                    // 最大空闲连接数
		sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(cfg.ConnMaxIdletime)) // 连接最大空闲时间
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)                                    // 最大打开连接数
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(cfg.ConnMaxLifetime)) // 连接最大生命周期

		sc.DBMap[cfg.DBName] = client
	}
}
