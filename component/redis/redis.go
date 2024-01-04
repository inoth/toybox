package redis

import (
	"github.com/redis/go-redis/v9"
)

type RedisComponent struct {
	name  string
	cache *redis.Client

	Addr        string `toml:"addr" json:"addr"`
	Username    string `toml:"username" json:"username"`
	Password    string `toml:"password" json:"password"`
	PoolSize    int    `toml:"pool_size" json:"pool_size"`
	PoolTimeout int    `toml:"pool_timeout" json:"pool_timeout"`
}

func NewRedisComponent(opts ...Option) *RedisComponent {
	return &RedisComponent{}
}
