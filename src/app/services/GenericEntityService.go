package services

import (
	"fmt"
	"rol/app/interfaces"
	"rol/app/interfaces/generic"
	"rol/app/mappers"
	"rol/app/utils"
	"rol/dtos"

	"github.com/sirupsen/logrus"
)

type GenericEntityService struct {
	repository generic.IGenericEntityRepository
	logger     *logrus.Logger
}

func NewGenericEntityService(repo generic.IGenericEntityRepository, log *logrus.Logger) (*GenericEntityService, error) {
	return &GenericEntityService{
		repository: repo,
		logger:     log,
	}, nil
}

func generateQuerySearchString(entity interface{}, search string) string {
	queryString := ""
	stringFieldNames := &[]string{}
	utils.GetStringFieldsNames(entity, stringFieldNames)
	snakeCaseColumnNames := toSnakeCase(stringFieldNames)
	for i := 0; i < len(*snakeCaseColumnNames); i++ {
		if i != 0 {
			queryString = queryString + " OR "
		}
		queryString = queryString + fmt.Sprintf("%s LIKE \"%s\"", (*snakeCaseColumnNames)[i], "%"+search+"%")
	}
	return queryString
}

func (ges *GenericEntityService) GetList(dtoArr interface{}, search, orderBy, orderDirection string, page, pageSize int) (*dtos.Paginator, error) {
	pageFinal := page
	pageSizeFinal := pageSize
	if page < 1 {
		pageFinal = 1
	}
	if pageSize < 1 {
		pageSizeFinal = 10
	}
	entities := mappers.GetEntityEmptyArray(dtoArr)

	queryRepString := ""
	var _ interface{} = nil
	if len(search) > 3 {
		entityModel := mappers.GetEmptyEntityFromArrayType(entities)
		queryRepString = generateQuerySearchString(entityModel, search)
	} else {
		queryRepString = ""
	}
	count, err := ges.repository.GetList(entities, orderBy, orderDirection, pageFinal, pageSizeFinal, queryRepString)
	if err != nil {
		ges.logger.Errorf("[Generic service] getting list error: %s", err.Error())
		return nil, err
	}

	mappers.Map(entities, dtoArr)
	return &dtos.Paginator{
		CurrentPage: pageFinal,
		PageSize:    pageSizeFinal,
		ItemsCount:  int(count),
		Items:       dtoArr,
	}, nil
}

func (ges *GenericEntityService) GetAll(dtoArr interface{}) error {
	entities := mappers.GetEntityEmptyArray(dtoArr)
	err := ges.repository.GetAll(entities)
	if err != nil {
		ges.logger.Errorf("[Generic service] error getting all entities: %s", err.Error())
		return err
	}
	mappers.Map(entities, dtoArr)
	return nil
}
func (ges *GenericEntityService) GetById(dto interfaces.IEntityDtoModel, id uint) error {
	entity := mappers.GetEmptyEntity(dto)
	err := ges.repository.GetById(entity, id)
	mappers.Map(entity, dto)
	if err != nil {
		ges.logger.Errorf("[Generic service] error retrieving an entity by its id: %s", err.Error())
		return err
	}
	return nil
}
func (ges *GenericEntityService) Update(updateDto interfaces.IEntityDtoModel, id uint) error {
	entity := mappers.GetEmptyEntity(updateDto)
	err := ges.repository.GetById(entity, id)
	if err != nil {
		ges.logger.Errorf("[Generic service] error retrieving an entity by its id: %s", err.Error())
		return err
	}
	mappers.Map(updateDto, entity)
	err = ges.repository.Update(entity)
	if err != nil {
		ges.logger.Errorf("[Generic service] error updating entity: %s", err.Error())
		return err
	}
	return nil
}
func (ges *GenericEntityService) Create(createDto interfaces.IEntityDtoModel) (uint, error) {
	entity := mappers.GetEmptyEntity(createDto)
	mappers.Map(createDto, entity)
	id, err := ges.repository.Insert(entity)
	if err != nil {
		ges.logger.Errorf("[Generic service] error creating entity: %s", err.Error())
		return 0, err
	}
	return id, nil
}
func (ges *GenericEntityService) Delete(dto interfaces.IEntityDtoModel, id uint) error {
	entity := mappers.GetEmptyEntity(dto)
	err := ges.repository.GetById(entity, id)
	if err != nil {
		ges.logger.Errorf("[Generic service] error retrieving an entity by its id: %s", err.Error())
		return err
	}
	err = ges.repository.Delete(entity)
	if err != nil {
		ges.logger.Errorf("[Generic service] error deleting entity: %s", err.Error())
	}
	return nil
}
