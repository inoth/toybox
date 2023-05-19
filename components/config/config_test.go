package config

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	cfg := New(func(ctx context.Context, cfg *ConfigComponent) {
		cfg.Interval = 3
		cfg.CfgPath = "../../../config/dev/global.toml"
	})
	err := cfg.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Printf("%+v", cfg.String())

	for {
		fmt.Println(Cfg.GetString("testa"))
		fmt.Println(Cfg.GetString("testb"))
		fmt.Println(Cfg.GetString("testc"))

		time.Sleep(time.Second * 5)
	}
}
