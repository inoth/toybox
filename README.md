# ino-toybox
项目手脚架, 常用组件封装, 模块式装载

~~稍微试试范型吧~~

```go
package main

import (
	"os"

	"github.com/gin-gonic/gin"
	inotoybox "github.com/inoth/ino-toybox"
	"github.com/inoth/ino-toybox/components/cache"
	"github.com/inoth/ino-toybox/components/config"
	"github.com/inoth/ino-toybox/components/logger"
	"github.com/inoth/ino-toybox/servers/httpsvc"
)

func main() {
	err := inotoybox.NewToyBox(
		&cache.CacheComponent{},
		config.Instance(),
		&logger.LogrusComponent{},
	).Init().Start(
		httpsvc.Instance(":8080").SetRouter(&DemoRouter{}),
	)
	if err != nil {
		os.Exit(1)
	}
}

type DemoRouter struct{}

func (DemoRouter) Prefix() string { return "" }
func (DemoRouter) LoadRouter(router *gin.RouterGroup) {
	router.GET("", func(c *gin.Context) {
		logger.Log.Info("ok")
		c.String(200, "ok")
	})
}
```