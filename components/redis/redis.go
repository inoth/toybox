package redis

import (
	"context"
	"time"

	gredis "github.com/go-redis/redis/v8"
	"github.com/inoth/toybox/components/config"
)

var Rdc *gredis.ClusterClient

// host:
//   - localhost:7000
//   - localhost:7001
//   - localhost:7002
//   - localhost:7003
//   - localhost:7004
//   - localhost:7005
// passwd: "1234567890"
// pool_size: 10
// pool_timeout: 3
type RedisComponent struct{}

func (RedisComponent) Init() error {

	hosts := config.Cfg.GetStringSlice("redis.host")
	password := config.Cfg.GetString("redis.passwd")
	poolSize := config.Cfg.GetInt("redis.pool_size")

	client := gredis.NewClusterClient(&gredis.ClusterOptions{
		Addrs:    hosts,
		Password: password,
		PoolSize: poolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*config.Cfg.GetDuration("redis.pool_timeout"))
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return err
	}
	Rdc = client
	return nil
}
