package rediscluster

import "github/inoth/toybox"

type RedisClusterComponent struct {
}

func New(tb *toybox.ToyBox) *RedisClusterComponent {
	return &RedisClusterComponent{}
}
