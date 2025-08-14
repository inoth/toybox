package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/inoth/toybox/config"
	"github.com/redis/go-redis/v9"
)

const (
	name = "redis"
)

type RedisComponent struct {
	rds []*redis.Client

	Addr        string `toml:"addr" json:"addr"`
	Username    string `toml:"username" json:"username"`
	Password    string `toml:"password" json:"password"`
	DB          []int  `toml:"db" json:"db"`
	PoolSize    int    `toml:"pool_size" json:"pool_size"`
	PoolTimeout int    `toml:"pool_timeout" json:"pool_timeout"`
}

func (rds *RedisComponent) Name() string {
	return name
}

func NewRedisCache(conf config.ConfigMate) *RedisComponent {
	rd := RedisComponent{}
	err := conf.PrimitiveDecode(&rd)
	if err != nil {
		panic(fmt.Errorf("init cache err: %v", err))
	}
	rd.initConnect()
	return &rd
}

func (rds *RedisComponent) initConnect() {
	for _, db := range rds.DB {
		client := redis.NewClient(&redis.Options{
			Addr:        rds.Addr,
			Username:    rds.Username,
			Password:    rds.Password,
			DB:          db,
			PoolSize:    rds.PoolSize,
			PoolTimeout: time.Duration(rds.PoolTimeout) * time.Second,
		})
		_, err := client.Ping(context.Background()).Result()
		if err != nil {
			panic(fmt.Errorf("failed to connect to redis: %v", err))
		}
		rds.rds = append(rds.rds, client)
	}
}

func (rds *RedisComponent) GetCache(idx ...int) *redis.Client {
	index := 0
	if len(idx) > 0 {
		index = idx[0]
	}
	return rds.rds[index]
}
