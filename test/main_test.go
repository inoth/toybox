package test

import (
	"fmt"
	"testing"

	inotoybox "github.com/inoth/ino-toybox"
	"github.com/inoth/ino-toybox/components/cache"
	"github.com/inoth/ino-toybox/components/config"
	"github.com/inoth/ino-toybox/components/logger"
	"github.com/inoth/ino-toybox/servers/wssvc"
)

func TestHttpSvc(t *testing.T) {

}

func TestChatSvc(t *testing.T) {
	err := inotoybox.NewToyBox(
		&cache.CacheComponent{},
		&config.ViperComponent{},
		&logger.LogrusComponent{},
	).Init().Start(&wssvc.HubServer{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
}
