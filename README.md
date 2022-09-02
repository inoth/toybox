# ino-toybox
**常用库使用初始化过程简化, 避免定义过多死板规则, 尽量还原原有操作**

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
	).Init().Start(
		httpsvc.Instance(":8080").
		SetRouter(&DemoRouter{}),
	)
	if err != nil {
		os.Exit(1)
	}
}

type DemoRouter struct{}

func (DemoRouter) Prefix() string { return "" }
func (DemoRouter) LoadRouter(router *gin.RouterGroup) {
	router.GET("", func(c *gin.Context) {
		c.String(200, "ok")
	})
}
```