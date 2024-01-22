package rediscluster

import "github.com/inoth/toybox"

var (
	default_name = "redis"
)

type Option func(*RedisClusterComponent)

func defaultOption() RedisClusterComponent {
	return RedisClusterComponent{
		name:        default_name,
		URLs:        make([]string, 0),
		Password:    "",
		PoolSize:    10,
		PoolTimeout: 60,
	}
}

func SetName(name string) Option {
	return func(mc *RedisClusterComponent) {
		mc.name = name
	}
}

func SetConfig(cfg toybox.ConfigMate) Option {
	return func(rc *RedisClusterComponent) {
		if err := cfg.PrimitiveDecodeComponent(rc); err != nil {
			panic("failed to load redis cluster configuration " + err.Error())
		}
	}
}
