package httpsvc

import (
	"github.com/gin-gonic/gin"
	"github.com/inoth/ino-toybox/components/config"
	"github.com/inoth/ino-toybox/services/httpsvc/middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GinRouters interface {
	Prefix() string
	LoadRouter(router *gin.RouterGroup)
}

type GinService struct {
	port   string
	engine *gin.Engine
}

func Instance(port ...string) *GinService {
	gsvc := &GinService{
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

func (gsvc *GinService) SetDefaultMiddleware() *GinService {
	gsvc.SetMiddleware(
		middleware.GinGlobalException(),
		middleware.Cors(),
	)
	return gsvc
}

func (gsvc *GinService) SetMiddleware(mids ...gin.HandlerFunc) *GinService {
	gsvc.engine.Use(mids...)
	return gsvc
}

func (gsvc *GinService) SetSwagger() *GinService {
	ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))
	gsvc.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return gsvc
}

func (gsvc *GinService) SetRouter(routers ...GinRouters) *GinService {
	for _, router := range routers {
		router.LoadRouter(gsvc.engine.Group(router.Prefix()))
	}
	return gsvc
}

func (gsvc *GinService) Start() error {
	return gsvc.engine.Run(gsvc.port)
}
