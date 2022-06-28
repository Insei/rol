package services

import (
	"fmt"
	"rol/app/interfaces"
	"rol/domain"
	"rol/dtos"

	"github.com/sirupsen/logrus"
)

//AppLogService service structure for AppLog entity
type AppLogService struct {
	*GenericService[dtos.AppLogDto,
		dtos.AppLogDto,
		dtos.AppLogDto,
		domain.AppLog]
}

//NewAppLogService preparing domain.AppLog repository for passing it in DI
//Params
//	rep - generic repository with http log instantiated
//	log - logrus logger
//Return
//	New app log service
func NewAppLogService(rep interfaces.IGenericRepository[domain.AppLog], log *logrus.Logger) (interfaces.IGenericService[
	dtos.AppLogDto,
	dtos.AppLogDto,
	dtos.AppLogDto,
	domain.AppLog], error) {
	genericService, err := NewGenericService[dtos.AppLogDto, dtos.AppLogDto, dtos.AppLogDto](rep, log)
	return &AppLogService{
		genericService,
	}, fmt.Errorf("error getting new log service: %s", err.Error())
}
