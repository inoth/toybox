package redis

import "github/inoth/toybox"

type RedisComponent struct {
}

func New(tb *toybox.ToyBox) *RedisComponent {
	return &RedisComponent{}
}
