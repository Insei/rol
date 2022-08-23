package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"reflect"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/mappers"
	"rol/app/utils"
	"rol/dtos"
	"strings"
)

//GenericService generic service structure for IEntityModel
type GenericService[DtoType interface{},
	CreateDtoType interface{},
	UpdateDtoType interface{},
	EntityType interfaces.IEntityModel] struct {
	repository    interfaces.IGenericRepository[EntityType]
	logger        *logrus.Logger
	logSourceName string
}

//NewGenericService generic service constructor
//Params
//	repo - generic repository
//	logger - logger
//Return
//	New generic service
func NewGenericService[DtoType interface{},
	CreateDtoType interface{},
	UpdateDtoType interface{},
	EntityType interfaces.IEntityModel](repo interfaces.IGenericRepository[EntityType], logger *logrus.Logger) (
	*GenericService[
		DtoType,
		CreateDtoType,
		UpdateDtoType,
		EntityType], error) {
	model := new(EntityType)
	return &GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]{
		repository:    repo,
		logger:        logger,
		logSourceName: fmt.Sprintf("GenericService<%s>", reflect.TypeOf(*model).Name()),
	}, nil
}

func (g *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) log(ctx context.Context, level, message string) {
	entry := g.logger.WithFields(logrus.Fields{
		"actionID": ctx.Value("requestID"),
		"source":   g.logSourceName,
	})
	switch level {
	case "err":
		entry.Error(message)
	case "info":
		entry.Info(message)
	case "warn":
		entry.Warn(message)
	case "debug":
		entry.Debug(message)
	}
}

func (g *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) excludeDeleted(queryBuilder interfaces.IQueryBuilder) {
	queryBuilder.Where("DeletedAt", "=", nil)
}

func (g *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) addSearchInAllFields(search string, queryBuilder interfaces.IQueryBuilder) {
	entityModel := new(EntityType)
	stringFieldNames := &[]string{}
	utils.GetStringFieldsNames(entityModel, stringFieldNames)
	queryGroup := g.repository.NewQueryBuilder(nil)
	for i := 0; i < len(*stringFieldNames); i++ {
		fieldName := (*stringFieldNames)[i]
		containPass := strings.Contains(strings.ToLower(fieldName), "pass")
		containKey := strings.Contains(strings.ToLower(fieldName), "key")
		if containPass || containKey {
			continue
		}

		queryGroup.Or(fieldName, "LIKE", "%"+search+"%")
	}
	queryBuilder = queryBuilder.WhereQuery(queryGroup)
}

func (g *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) getListBasic(ctx context.Context, queryBuilder interfaces.IQueryBuilder, orderBy, orderDirection string, page, pageSize int) (*dtos.PaginatedItemsDto[DtoType], error) {
	pageFinal := page
	pageSizeFinal := pageSize
	if page < 1 {
		pageFinal = 1
	}
	if pageSize < 1 {
		pageSizeFinal = 10
	}
	entities, err := g.repository.GetList(ctx, orderBy, orderDirection, pageFinal, pageSizeFinal, queryBuilder)
	if err != nil {
		return nil, errors.Internal.New("repository failed get list")
	}
	count, err := g.repository.Count(ctx, queryBuilder)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error counting entities")
	}
	dtoArr := &[]DtoType{}
	for i := 0; i < len(*entities); i++ {
		dto := new(DtoType)
		err = mappers.MapEntityToDto((*entities)[i], dto)
		if err != nil {
			return nil, errors.Internal.Wrap(err, "error map entity to dto")
		}
		*dtoArr = append(*dtoArr, *dto)
	}

	paginatedDto := new(dtos.PaginatedItemsDto[DtoType])
	paginatedDto.Pagination.Page = pageFinal
	paginatedDto.Pagination.Size = pageSizeFinal
	paginatedDto.Pagination.TotalCount = int(count)
	paginatedDto.Pagination.TotalPages = pageSizeFinal / pageFinal
	if pageSizeFinal%pageFinal > 0 {
		paginatedDto.Pagination.TotalPages++
	}
	paginatedDto.Items = *dtoArr
	return paginatedDto, nil
}

//GetList Get list of elements with filtering and pagination
//Params
//	ctx - context is used only for logging
//	search - string for search in entity string fields
//	orderBy - order by entity field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return
//	*dtos.PaginatedItemsDto[DtoType] - pointer to paginated list
//	error - if an error occurs, otherwise nil
func (g *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GetList(ctx context.Context, search, orderBy, orderDirection string, page, pageSize int) (*dtos.PaginatedItemsDto[DtoType], error) {
	searchQueryBuilder := g.repository.NewQueryBuilder(ctx)
	g.excludeDeleted(searchQueryBuilder)
	if len(search) > 3 {
		g.addSearchInAllFields(search, searchQueryBuilder)
	}
	list, err := g.getListBasic(ctx, searchQueryBuilder, orderBy, orderDirection, page, pageSize)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to get list ")
	}
	return list, nil
}

func (g *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) getByIDBasic(ctx context.Context, id uuid.UUID, queryBuilder interfaces.IQueryBuilder) (*DtoType, error) {
	entity, err := g.repository.GetByIDExtended(ctx, id, queryBuilder)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to get by id")
	}
	if entity == nil {
		return nil, nil
	}
	dto := new(DtoType)
	err = mappers.MapEntityToDto(*entity, dto)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error map entity to dto")
	}
	return dto, nil
}

//GetByID Get entity by ID
//Params
//	ctx - context is used only for logging
//	id - entity id
//Return
//	*DtoType - pointer to dto
//	error - if an error occurs, otherwise nil
func (g *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GetByID(ctx context.Context, id uuid.UUID) (*DtoType, error) {
	queryBuilder := g.repository.NewQueryBuilder(ctx)
	g.excludeDeleted(queryBuilder)
	return g.getByIDBasic(ctx, id, queryBuilder)
}

func (g *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) updateBasic(ctx context.Context, updateDto UpdateDtoType, id uuid.UUID, queryBuilder interfaces.IQueryBuilder) (DtoType, error) {
	entity, err := g.repository.GetByIDExtended(ctx, id, queryBuilder)
	if err != nil {
		return *new(DtoType), errors.Internal.Wrap(err, "failed to get entity by id")
	}
	if entity == nil {
		return *new(DtoType), errors.NotFound.Wrap(err, "entity not found")
	}
	err = mappers.MapDtoToEntity(updateDto, entity)
	if err != nil {
		return *new(DtoType), errors.Internal.Wrap(err, "error map dto to entity")
	}
	updEntity, err := g.repository.Update(ctx, entity)
	if err != nil {
		return *new(DtoType), errors.Internal.Wrap(err, "failed to update entity in repository")
	}
	dto := new(DtoType)
	err = mappers.MapEntityToDto(updEntity, dto)
	if err != nil {
		return *new(DtoType), errors.Internal.Wrap(err, "error map entity to dto")
	}
	return *dto, nil
}

//Update save the changes to the existing entity
//Params
//	ctx - context is used only for logging
//	updateDto - update dto
//	id - entity id
//Return
//	DtoType - updated entity dto
//	error - if an error occurs, otherwise nil
func (g *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) Update(ctx context.Context, updateDto UpdateDtoType, id uuid.UUID) (DtoType, error) {
	queryBuilder := g.repository.NewQueryBuilder(ctx)
	g.excludeDeleted(queryBuilder)
	return g.updateBasic(ctx, updateDto, id, queryBuilder)
}

//Create add new entity
//Params
//	ctx - context is used only for logging
//	createDto - create dto
//Return
//	DtoType - created entity dto
//	error - if an error occurs, otherwise nil
func (g *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) Create(ctx context.Context, createDto CreateDtoType) (DtoType, error) {
	entity := new(EntityType)
	outDto := new(DtoType)
	err := mappers.MapDtoToEntity(createDto, entity)
	if err != nil {
		return *outDto, errors.Internal.Wrap(err, "error map entity to dto")
	}
	newEntity, err := g.repository.Insert(ctx, *entity)
	if err != nil {
		return *outDto, errors.Internal.Wrap(err, "create entity error")
	}
	err = mappers.MapEntityToDto(newEntity, outDto)
	if err != nil {
		return *outDto, errors.Internal.Wrap(err, "error map dto to entity")
	}
	return *outDto, nil
}

//Delete mark entity as deleted
//Params
//	ctx - context is used only for logging
//	id - entity id
//Return
//	error - if an error occurs, otherwise nil
func (g *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) Delete(ctx context.Context, id uuid.UUID) error {
	queryBuilder := g.repository.NewQueryBuilder(ctx)
	g.excludeDeleted(queryBuilder)
	entity, err := g.repository.GetByIDExtended(ctx, id, queryBuilder)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to get entity by id")
	}
	if entity == nil {
		return errors.NotFound.Wrap(err, "entity not found")
	}

	//entityReflect := reflect.ValueOf(entity).Interface().(interfaces.IEntityModelDeletedAt)
	//entityReflect.SetDeleted()

	err = g.repository.Delete(ctx, id)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to update entity in the repository")
	}
	return nil
}
