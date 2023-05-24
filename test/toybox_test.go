package test

import (
	"os"
	"testing"

	"github.com/inoth/toybox"
	"github.com/inoth/toybox/components/config"
	"github.com/inoth/toybox/components/mysql"
	"github.com/inoth/toybox/components/redis"
)

func TestNewToyBox(t *testing.T) {
	os.Setenv("GORUNEVN", "dev")

	tb := toybox.New(
		toybox.WithComponentCfgPath("../config"),
		toybox.EnableComponents(config.New(), redis.New(), mysql.New()),
	)

	tb.Run()
}
