package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	once sync.Once
	rdc  *RedisComponent
)

type RedisComponent struct {
	name  string
	ready bool
	cache *redis.Client

	Addr        string `toml:"addr" json:"addr"`
	Username    string `toml:"username" json:"username"`
	Password    string `toml:"password" json:"password"`
	DB          int    `toml:"db" json:"db"`
	PoolSize    int    `toml:"pool_size" json:"pool_size"`
	PoolTimeout int    `toml:"pool_timeout" json:"pool_timeout"`
}

func (rds RedisComponent) Name() string {
	return rds.name
}

func (rds RedisComponent) Ready() bool {
	return rds.ready
}

func (rds *RedisComponent) IsReady() {
	rds.ready = true
}

func (rds *RedisComponent) Init(ctx context.Context) error {
	return nil
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
	rds.cache = client
}

func GetCache(opts ...Option) (*redis.Client, error) {
	var err error
	once.Do(func() {
		rc := defaultOption()
		for _, opt := range opts {
			opt(&rc)
		}
		if !rc.Ready() {
			err = fmt.Errorf("components %s not yet ready", rc.name)
			return
		}
		rc.newCache()
		rdc = &rc
	})
	if err != nil {
		return nil, err
	}
	return rdc.cache, nil
}
