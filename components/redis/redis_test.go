package redis

import (
	"fmt"
	"testing"
)

func TestConnectToRedis(t *testing.T) {
	var (
		URLs   = []string{""}
		passwd = ""
	)
	client := &RedisComponents{
		URLs:     URLs,
		Password: passwd,
	}
	client.Set("testkey", 123)
	res, _ := client.Get("testkey")
	fmt.Println(res)
}
