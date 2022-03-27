package services

import (
	// "fmt"
	// "regexp"
	"errors"
	"fmt"
	"reflect"
	"rol/app/interfaces"
	"rol/app/interfaces/generic"
	"rol/app/mappers"
	"rol/app/utils"
	"rol/dtos"
	"strings"
	"time"

	// "rol/app/mappers"
	// "rol/app/utils"
	// "rol/dtos"
	// "strings"

	"github.com/sirupsen/logrus"
)

//TODO: remove logrus and switch to the logger interface
type GenericService[DtoType interfaces.IEntityDtoModel,
	CreateDtoType interfaces.IEntityDtoModel,
	UpdateDtoType interfaces.IEntityDtoModel,
	EntityType interfaces.IEntityModel] struct {
	repository generic.IGenericRepository[EntityType]
	logger     *logrus.Logger
}

func NewGenericService[DtoType interfaces.IEntityDtoModel,
	CreateDtoType interfaces.IEntityDtoModel,
	UpdateDtoType interfaces.IEntityDtoModel,
	EntityType interfaces.IEntityModel](repo generic.IGenericRepository[EntityType], logger *logrus.Logger) (
	*GenericService[
		DtoType,
		CreateDtoType,
		UpdateDtoType,
		EntityType], error) {
	return &GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]{
		repository: repo,
		logger:     logger,
	}, nil
}

func (s *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) getEntityTypeName() string {
	model := new(EntityType)
	return reflect.TypeOf(*model).Name()
}

func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) excludeDelted(queryBuilder interfaces.IQueryBuilder) {
	date, err := time.Parse("2006-01-02", "1999-01-01")
	if err != nil {
		ges.logger.Errorf("time parse error: %s", err.Error())
		return
	}
	queryBuilder.Where("DeletedAt", "<", date)
}

func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) addSearchInAllFields(search string, queryBuilder interfaces.IQueryBuilder) {
	entityModel := new(EntityType)
	stringFieldNames := &[]string{}
	utils.GetStringFieldsNames(entityModel, stringFieldNames)
	for i := 0; i < len(*stringFieldNames); i++ {
		fieldName := (*stringFieldNames)[i]
		containPass := strings.Contains(strings.ToLower(fieldName), "pass")
		containKey := strings.Contains(strings.ToLower(fieldName), "key")
		if containPass || containKey {
			continue
		}
		queryBuilder = queryBuilder.Where(fieldName, "LIKE", "%"+search+"%")
	}
}

func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GetList(search, orderBy, orderDirection string, page, pageSize int) (*dtos.PaginatedListDto[DtoType], error) {
	pageFinal := page
	pageSizeFinal := pageSize
	if page < 1 {
		pageFinal = 1
	}
	if pageSize < 1 {
		pageSizeFinal = 10
	}
	searchQueryBuilder := ges.repository.NewQueryBuilder()
	ges.excludeDelted(searchQueryBuilder)
	if len(search) > 3 {
		ges.addSearchInAllFields(search, searchQueryBuilder)
	}
	dtoArr := new([]DtoType)
	entities, err := ges.repository.GetList(orderBy, orderDirection, pageFinal, pageSizeFinal, searchQueryBuilder)
	if err != nil {
		return nil, err
	}
	count, err := ges.repository.Count(searchQueryBuilder)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(*entities); i++ {
		entityDto := new(DtoType)
		mappers.Map((*entities)[i], entityDto)
		*dtoArr = append(*dtoArr, *entityDto)
	}

	paginatedDto := new(dtos.PaginatedListDto[DtoType])
	paginatedDto.Page = pageFinal
	paginatedDto.Size = pageSizeFinal
	paginatedDto.Total = count
	paginatedDto.Items = dtoArr
	return paginatedDto, nil
}

func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) getByIdExcludeDelted(id uint) (*EntityType, error) {
	queryBuilder := ges.repository.NewQueryBuilder()
	ges.excludeDelted(queryBuilder)
	queryBuilder.Where("id", "=", id)
	entities, err := ges.repository.GetList("id", "DESC", 1, 1, queryBuilder)
	if err != nil {
		return nil, err
	}
	if len(*entities) == 0 {
		return nil, nil
	}
	entity := &(*entities)[0]
	return entity, nil
}

func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GetById(id uint) (*DtoType, error) {
	entity, err := ges.getByIdExcludeDelted(id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}
	entityDto := new(DtoType)
	mappers.Map(entity, entityDto)
	return entityDto, nil
}

func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) Update(updateDto UpdateDtoType, id uint) error {
	entity, err := ges.getByIdExcludeDelted(id)
	if err != nil {
		return err
	}
	if entity == nil {
		errStr := fmt.Sprintf("[Generic Service]: [%s]: [UPDATE]: entity with id %d does not exist", ges.getEntityTypeName(), id)
		err = errors.New(errStr)
		return err
	}
	mappers.Map(updateDto, entity)
	err = ges.repository.Update(entity)
	if err != nil {
		return err
	}
	return nil
}

func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) Create(createDto CreateDtoType) (uint, error) {
	entity := new(EntityType)
	mappers.Map(createDto, entity)
	id, err := ges.repository.Insert(*entity)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) Delete(id uint) error {
	entity, err := ges.getByIdExcludeDelted(id)
	if err != nil {
		return err
	}
	if entity == nil {
		errStr := fmt.Sprintf("[Generic Service]: [%s]: [DELETE]: entity with id %d does not exist", ges.getEntityTypeName(), id)
		err = errors.New(errStr)
		return err
	}
	(*entity).SetDeleted()
	err = ges.repository.Update(entity)
	if err != nil {
		return err
	}
	return nil
}
