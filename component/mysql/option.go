package mysql

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

func WithDbName(name string) Option {
	return func(mc *MysqlComponent) {
		mc.name = name
	}
}
