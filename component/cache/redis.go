package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/inoth/toybox/config"
	"github.com/redis/go-redis/v9"
)

const (
	Name = "redis"
)

type RedisComponent struct {
	rd *redis.Client

	Addr        string `toml:"addr" json:"addr"`
	Username    string `toml:"username" json:"username"`
	Password    string `toml:"password" json:"password"`
	DB          int    `toml:"db" json:"db"`
	PoolSize    int    `toml:"pool_size" json:"pool_size"`
	PoolTimeout int    `toml:"pool_timeout" json:"pool_timeout"`
}

func (rds *RedisComponent) Name() string {
	return Name
}

func NewCache(conf config.ConfigMate) *RedisComponent {
	rd := RedisComponent{}
	err := conf.PrimitiveDecode(&rd)
	if err != nil {
		panic(fmt.Errorf("init mysql err: %v", err))
	}
	rd.newCache()
	return &rd
}

func (rds *RedisComponent) newCache() {
	client := redis.NewClient(&redis.Options{
		Addr:        rds.Addr,
		Username:    rds.Username,
		Password:    rds.Password,
		DB:          rds.DB,
		PoolSize:    rds.PoolSize,
		PoolTimeout: time.Duration(rds.PoolTimeout) * time.Second,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Errorf("failed to connect to redis: %v", err))
	}
	rds.rd = client
}

func (rds *RedisComponent) GetCache() *redis.Client {
	return rds.rd
}
