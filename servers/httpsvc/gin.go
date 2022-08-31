package httpsvc

import (
	"github.com/gin-gonic/gin"
	"github.com/inoth/ino-toybox/components/config"
	"github.com/inoth/ino-toybox/servers/httpsvc/middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GinRouters interface {
	Prefix() string
	LoadRouter(router *gin.RouterGroup)
}

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
		gsvc.port = config.Cfg.GetString("ServerPort")
	}
	gsvc.engine.MaxMultipartMemory = 10 << 20
	return gsvc
}

func (gsvc *GinServer) SetDefaultMiddleware() *GinServer {
	gsvc.SetMiddleware(
		middleware.GinGlobalException(),
		middleware.Cors(),
	)
	return gsvc
}

func (gsvc *GinServer) SetMiddleware(mids ...gin.HandlerFunc) *GinServer {
	gsvc.engine.Use(mids...)
	return gsvc
}

func (gsvc *GinServer) SetSwagger() *GinServer {
	ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))
	gsvc.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return gsvc
}

func (gsvc *GinServer) SetRouter(routers ...GinRouters) *GinServer {
	for _, router := range routers {
		router.LoadRouter(gsvc.engine.Group(router.Prefix()))
	}
	return gsvc
}

func (gsvc *GinServer) Start() error {
	return gsvc.engine.Run(gsvc.port)
}
