package mysql

import "github/inoth/toybox/conf"

type Option func(*MysqlComponent)

func defaultOption(name string) MysqlComponent {
	if name == "" {
		name = default_name
	}
	return MysqlComponent{
		name:            name,
		Host:            "localhost",
		Port:            3306,
		User:            "root",
		Passwd:          "",
		DbName:          name,
		MaxIdleConns:    100,
		MaxOpenConns:    100,
		ConnMaxIdletime: 60,
		ConnMaxLifetime: 60,
	}
}

func SetConfig(cfg conf.Conf) Option {
	return func(mc *MysqlComponent) {
		err := cfg.Configuration(mc.name, mc)
		if err != nil {
			panic("failed to load configuration")
		}
		mc.ready = true
	}
}
