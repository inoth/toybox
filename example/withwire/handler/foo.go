package handler

import (
	"demo/service"
	"net/http"

	"github.com/inoth/toybox/component/logger"
	"github.com/inoth/toybox/server/ginsvr"

	"github.com/gin-gonic/gin"
)

type FooHandler struct {
	log *logger.Logger
	svc *service.FooService

	prefix  string
	mid     []gin.HandlerFunc
	routers []ginsvr.Router
}

func NewFooHandler(svc *service.FooService) *FooHandler {
	foo := &FooHandler{
		prefix: "api",
		svc:    svc,
		log:    logger.GetLogger(logger.LoggerConfig{ServerName: "FooController"}),
	}
	foo.AppendRouter(http.MethodGet, "hi", foo.SayHi())
	return foo
}

func (foo *FooHandler) AppendRouter(method, path string, handles ...gin.HandlerFunc) {
	foo.routers = append(foo.routers, ginsvr.NewRouter(method, path, handles...))
}

func (foo FooHandler) Prefix() string {
	return foo.prefix
}

func (foo FooHandler) Routers() []ginsvr.Router {
	return foo.routers
}

func (foo FooHandler) Middlewares() []gin.HandlerFunc {
	return foo.mid
}

func (foo *FooHandler) SayHi() gin.HandlerFunc {
	return func(c *gin.Context) {
		foo.log.Info(foo.svc.SayHi())
		c.String(200, foo.svc.SayHi())
	}
}
