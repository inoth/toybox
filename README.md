# ino-toybox

```go
package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/register"
	"github.com/inoth/toybox/components/cache"
	"github.com/inoth/toybox/components/config"
	"github.com/inoth/toybox/components/logger"
	"github.com/inoth/toybox/servers/httpsvc"
)

func main() {
	err := register.NewToyBox(
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