package redis

import (
	"fmt"

	"github.com/inoth/toybox"
)

var (
	default_name = "redis"
)

type Option func(*RedisComponent)

func defaultOption() RedisComponent {
	return RedisComponent{
		name:        default_name,
		Addr:        "localhost:6379",
		Password:    "",
		PoolSize:    10,
		PoolTimeout: 60,
	}
}

func SetName(name string) Option {
	return func(mc *RedisComponent) {
		mc.name = name
	}
}

func SetConfig(cfg toybox.ConfigMate) Option {
	return func(rc *RedisComponent) {
		if err := cfg.PrimitiveDecodeComponent(rc); err != nil {
			panic(fmt.Errorf("failed to load redis configuration; %v", err))
		}
	}
}
