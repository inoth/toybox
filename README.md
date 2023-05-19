# toybox

通过 go get 获取包
```shell
go get -u github.com/inoth/toybox
```
---

main.go
```go
package main

import (
	"os"
	"testing"

	"github.com/inoth/toybox"
	"github.com/inoth/toybox/components/config"
	"github.com/inoth/toybox/components/mysql"
	"github.com/inoth/toybox/components/redis"

	_ "webv1"
)

func main() {
	tb := toybox.New(
		toybox.WithCfgPath("config"),
		toybox.EnableComponents(
			config.New(), 
			redis.New(), 
			mysql.New(),
		),
	)

	tb.Run()
}
```
webapi.go
```go
package webv1

import (
	"github.com/inoth/toybox/server"
	"github.com/gin-gonic/gin"
)

var (
	serverName = "web"
	ver        = "v1"
)

type ChatServer struct {
	ServerName string
	Port       string
	// 依赖组件
	RequiredComponents []string
}

func init() {
	server.Add(serverName, func() server.Server {
		return &ChatServer{
			ServerName:         serverName,
			Port:               ":8888",
			RequiredComponents: []string{"config", "redis", "mysql"},
		}
	})
}

func (cs *ChatServer) Name() string {return serverName}

func (cs *ChatServer) RequiredComponent() []string {
	return cs.RequiredComponents
}

func (cs *ChatServer) Start() error {
	r := gin.New()
	chat := r.Group(cs.ServerName + "/" + ver)
	{ // heartbeat
		chat.GET("heartbeat", func(ctx *gin.Context) { 
			ctx.String(200, "")
		})
	}
	{ // business router
		
	}
	r.Run(cs.Port)
	return nil
}
```
