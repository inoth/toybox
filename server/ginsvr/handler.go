package ginsvr

import "github.com/gin-gonic/gin"

type Router struct {
	Method string
	Path   string
	Handle []gin.HandlerFunc
}

func NewRouter(method, path string, handles ...gin.HandlerFunc) Router {
	return Router{
		Method: method,
		Path:   path,
		Handle: handles,
	}
}

type Handler interface {
	Prefix() string
	Middlewares() []gin.HandlerFunc
	Routers() []Router
}
