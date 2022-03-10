package webapi

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rol/app/interfaces/generic"
	"rol/webapi/controllers"
	"rol/webapi/middleware"
)

type HttpServer struct {
	engine  *gin.Engine
	service *generic.IGenericEntityService
	logger  *logrus.Logger
}

func NewHttpServer(log *logrus.Logger, service *generic.IGenericEntityService) HttpServer {
	ginEngine := gin.New()
	ginEngine.Use(middleware.Logger(log), middleware.Recovery(log))
	//ginEngine.Use(middleware.GinBodyLogMiddleware)
	server := HttpServer{
		engine:  ginEngine,
		service: service,
		logger:  log,
	}
	return server
}

func (server *HttpServer) InitializeRoutes() {
	server.InitializeControllers()
}

func (server *HttpServer) InitializeControllers() {
	switchContr := controllers.NewEthernetSwitchController(server.service)

	groupRoute := server.engine.Group("/api/v1")

	groupRoute.GET("/switch/list", switchContr.GetList)
	groupRoute.GET("/switch/:id", switchContr.GetById)
	groupRoute.GET("/switch", switchContr.GetAll)
	groupRoute.POST("/switch", switchContr.Create)
	groupRoute.PUT("/switch/:id", switchContr.Update)
}

func (server *HttpServer) Start(address string) {
	server.InitializeRoutes()
	err := server.engine.Run(address)
	if err != nil {
		server.logger.Errorf("[Http server] start server error: %s", err.Error())
		return
	}
}
