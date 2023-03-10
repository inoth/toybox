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

func LoadRouter(router GinRouters) {
	if _, ok := RoutersMap[router.Prefix()]; !ok {
		RoutersMap[router.Prefix()] = router
	}
	fmt.Println("重复添加服务路由: ", router.Prefix())
}
