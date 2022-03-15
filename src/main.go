package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"os"
	"rol/app"
	"rol/app/interfaces/generic"
	"rol/app/services"
	"rol/infrastructure"
	"rol/webapi"
	"rol/webapi/utils"
)

func main() {
	cfg := app.GetConfig()
	// We need use service as interface, and not as the struct, then we can see implementation errors.
	var service generic.IGenericEntityService = nil
	var repository generic.IGenericEntityRepository = nil
	// Setup sql connection
	gormSqlConnection := mysql.Open(cfg.Database.SqlConnectionString)
	// Setup logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&utils.Formatter{})
	// Setup generic repo (infrastructure layer)
	repository, _ = infrastructure.NewGormGenericEntityRepository(gormSqlConnection, logger)
	//Setup Generic service (business layer, i.e. app)
	service, _ = services.NewGenericEntityService(repository, logger)
	// Setup http server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	httpServer := webapi.NewHttpServer(logger, &service)
	httpServer.Start(addr)

}
