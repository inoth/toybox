package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/inoth/toybox/component"

	gredis "github.com/go-redis/redis/v8"
)

var (
	componentName = "redis"
	redisOnce     sync.Once
	Rdc           *RedisComponents
)

// func init() {
// 	components.Add(componentName, func() toybox.Component {
// 		return &RedisComponents{}
// 	})
// }

func New() component.Component {
	return &RedisComponents{}
}

type RedisComponents struct {
	client *gredis.ClusterClient

	URLs        []string `toml:"urls" json:"urls" yaml:"urls"`
	Username    string   `toml:"username" json:"username" yaml:"username"`
	Password    string   `toml:"password" json:"password" yaml:"password"`
	PoolSize    int      `toml:"pool_size" json:"pool_size" yaml:"pool_size"`
	PoolTimeout int      `toml:"pool_time_out" json:"pool_time_out" yaml:"pool_time_out"`
}

func (rc *RedisComponents) Name() string {
	return componentName
}

func (rc *RedisComponents) String() string {
	buf, _ := json.Marshal(rc)
	return string(buf)
}

func (rc *RedisComponents) Close() error { return rc.client.Close() }

func (rc *RedisComponents) Init() (err error) {
	redisOnce.Do(func() {
		client := gredis.NewClusterClient(&gredis.ClusterOptions{
			Addrs:       rc.URLs,
			Password:    rc.Password,
			Username:    rc.Username,
			PoolSize:    rc.PoolSize,
			PoolTimeout: time.Duration(rc.PoolTimeout),
		})
		_, err = client.Ping(context.Background()).Result()
		if err != nil {
			err = fmt.Errorf("failed to connect to redis: %v", err)
			return
		}
		rc.client = client
		Rdc = rc
		fmt.Println("redis component initialization successful")
	})
	return
}

func (rc *RedisComponents) Get(key string) (string, error) {
	res, err := rc.client.Get(context.Background(), key).Result()
	if err == gredis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return res, nil
}

func (rc *RedisComponents) Set(key string, val interface{}, expiration ...time.Duration) error {
	expir := time.Duration(0)
	if len(expiration) > 0 {
		expir = expiration[0]
	}
	return rc.client.Set(context.Background(), key, val, expir).Err()
}

func (rc *RedisComponents) Redis() *gredis.ClusterClient {
	return rc.client
}
