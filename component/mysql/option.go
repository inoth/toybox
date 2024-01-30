package mysql

import (
	"fmt"

	"github.com/inoth/toybox"
)

const (
	default_name = "mysql"
)

type Option func(*MysqlComponent)

func defaultOption() MysqlComponent {
	return MysqlComponent{
		name:            default_name,
		Host:            "localhost",
		Port:            3306,
		User:            "root",
		Passwd:          "",
		DbName:          "mysql",
		MaxIdleConns:    100,
		MaxOpenConns:    100,
		ConnMaxIdletime: 60,
		ConnMaxLifetime: 60,
	}
}

func SetName(name string) Option {
	return func(mc *MysqlComponent) {
		mc.name = name
	}
}

func SetConfig(cfg toybox.ConfigMate) Option {
	return func(mc *MysqlComponent) {
		if err := cfg.PrimitiveDecodeComponent(mc); err != nil {
			panic(fmt.Errorf("failed to load mysql configuration; %v", err))
		}
	}
}
