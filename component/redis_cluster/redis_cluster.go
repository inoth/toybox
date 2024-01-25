package rediscluster

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	once sync.Once
	rdc  *RedisClusterComponent
)

type RedisClusterComponent struct {
	name  string
	ready bool
	cache *redis.ClusterClient

	URLs        []string `toml:"urls" json:"urls"`
	Username    string   `toml:"username" json:"username"`
	Password    string   `toml:"password" json:"password"`
	PoolSize    int      `toml:"pool_size" json:"pool_size"`
	PoolTimeout int      `toml:"pool_timeout" json:"pool_timeout"`
}

func (rds *RedisClusterComponent) Name() string {
	return rds.name
}

func (rds *RedisClusterComponent) Ready() bool {
	return rds.ready
}

func (rds *RedisClusterComponent) IsReady() {
	rds.ready = true
}

func (rds *RedisClusterComponent) Init(ctx context.Context) error {
	return nil
}

func (rds *RedisClusterComponent) newCache() {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:       rds.URLs,
		Password:    rds.Password,
		Username:    rds.Username,
		PoolSize:    rds.PoolSize,
		PoolTimeout: time.Duration(rds.PoolTimeout) * time.Second,
	})
	err := client.ForEachShard(context.Background(), func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		panic(fmt.Errorf("failed to connect to redis cluster: %v", err))
	}
	rds.cache = client
}

func GetCache(opts ...Option) (*redis.ClusterClient, error) {
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
