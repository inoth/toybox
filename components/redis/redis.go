package redis

import (
	"context"
	"time"

	gredis "github.com/go-redis/redis/v8"
	"github.com/inoth/ino-toybox/components/config"
)

var Rdc *gredis.ClusterClient

// Redis:
//   Host:
//     - localhost:7000
//     - localhost:7001
//     - localhost:7002
//     - localhost:7003
//     - localhost:7004
//     - localhost:7005
//   Passwd: "123456789"
//   PoolSize: 10
//   PoolTimeout: 3
type RedisConnectCluster struct{}

func (RedisConnectCluster) Init() error {

	hosts := config.Cfg.GetStringSlice("Redis.Host")
	password := config.Cfg.GetString("Redis.Passwd")
	poolSize := config.Cfg.GetInt("Redis.PoolSize")

	client := gredis.NewClusterClient(&gredis.ClusterOptions{
		Addrs:    hosts,
		Password: password,
		PoolSize: poolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*config.Cfg.GetDuration("Redis.PoolTimeout"))
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return err
	}
	Rdc = client
	return nil
}
