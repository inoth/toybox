package ginsvc

import (
	"github.com/gin-gonic/gin"
)

type GinServer struct {
	port   string
	engine *gin.Engine
}

func Instance(port ...string) *GinServer {
	gsvc := &GinServer{
		engine: gin.New(),
	}
	if len(port) > 0 {
		gsvc.port = port[0]
	} else {
		gsvc.port = ":8080"
	}
	gsvc.engine.MaxMultipartMemory = 10 << 20
	return gsvc
}

func (gsvc *GinServer) SetMiddleware(mids ...gin.HandlerFunc) *GinServer {
	gsvc.engine.Use(mids...)
	return gsvc
}

func (gsvc *GinServer) SetRouter(routers ...GinRouters) *GinServer {
	for _, router := range routers {
		router.LoadRouter(gsvc.engine.Group(router.Prefix()))
	}
	for _, router := range RoutersMap {
		router.LoadRouter(gsvc.engine.Group(router.Prefix()))
	}
	return gsvc
}

func (gsvc *GinServer) Start() error {
	return gsvc.engine.Run(gsvc.port)
}

func (gsvc *GinServer) Stop() {}
