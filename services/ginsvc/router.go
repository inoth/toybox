package ginsvc

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type GinRouters interface {
	Prefix() string
	LoadRouter(router *gin.RouterGroup)
}

var RoutersMap = map[string]GinRouters{}

func LoadRouter(serverName string, router GinRouters) {
	if _, ok := RoutersMap[serverName]; !ok {
		RoutersMap[serverName] = router
	}
	fmt.Println("重复添加服务路由: ", serverName)
}
