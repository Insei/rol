package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/mappers"
	"rol/domain"
	"rol/dtos"
)

//UserProjectService service structure for user project entity
type UserProjectService struct {
	subnets    map[string]bool
	repository interfaces.IGenericRepository[domain.Project]
	manager    interfaces.IProjectManager
	logger     *logrus.Logger
}

//NewUserProjectService constructor for UserProjectService
func NewUserProjectService(repo interfaces.IGenericRepository[domain.Project], manager interfaces.IProjectManager, log *logrus.Logger) (*UserProjectService, error) {
	subnets := make(map[string]bool)
	subnets["10.10.10.0"] = true
	subnets["10.10.11.0"] = true
	subnets["10.10.12.0"] = true
	subnets["10.10.13.0"] = true
	subnets["10.10.14.0"] = true
	subnets["10.10.15.0"] = true

	return &UserProjectService{
		subnets:    subnets,
		repository: repo,
		logger:     log,
		manager:    manager,
	}, nil
}

//GetList Get list of user projects dto elements with search in all fields and pagination
//
//Params:
//	ctx - context is used only for logging
//	search - string for search in entity string fields
//	orderBy - order by entity field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return:
//	dtos.PaginatedItemsDto[dtos.ProjectDto] - paginated items
//	error - if an error occurs, otherwise nil
func (p *UserProjectService) GetList(ctx context.Context, search string, orderBy string, orderDirection string, page int, pageSize int) (dtos.PaginatedItemsDto[dtos.ProjectDto], error) {
	return GetList[dtos.ProjectDto](ctx, p.repository, search, orderBy, orderDirection, page, pageSize)
}

//GetByID Get user project dto by ID
//
//Params:
//	ctx - context is used only for logging
//	id - user project id
//Return:
//	dtos.ProjectDto - user project dto
//	error - if an error occurs, otherwise nil
func (p *UserProjectService) GetByID(ctx context.Context, id uuid.UUID) (dtos.ProjectDto, error) {
	return GetByID[dtos.ProjectDto](ctx, p.repository, id, nil)
}

//Create add new user project
//Params
//	ctx - context
//	createDto - user project create dto
//Return
//	dtos.ProjectDto - created user project
//	error - if an error occurs, otherwise nil
func (p *UserProjectService) Create(ctx context.Context, createDto dtos.ProjectCreateDto) (dtos.ProjectDto, error) {
	entity := new(domain.Project)
	outDto := new(dtos.ProjectDto)
	err := mappers.MapDtoToEntity(createDto, entity)
	if err != nil {
		return *outDto, errors.Internal.Wrap(err, "error map entity to dto")
	}
	project, err := p.manager.CreateProject(*entity, createDto.Subnet)
	if err != nil {
		return dtos.ProjectDto{}, errors.Internal.Wrap(err, "project manager failed to create project")
	}
	newEntity, err := p.repository.Insert(ctx, project)
	if err != nil {
		return *outDto, errors.Internal.Wrap(err, "repository failed to insert project")
	}
	err = mappers.MapEntityToDto(newEntity, outDto)
	if err != nil {
		return *outDto, errors.Internal.Wrap(err, "error map dto to entity")
	}
	return *outDto, nil
}

//Delete mark user project as deleted
//Params
//	ctx - context is used only for logging
//	id - user project id
//Return
//	error - if an error occurs, otherwise nil
func (p *UserProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	project, err := p.repository.GetByID(ctx, id)
	if err != nil {
		return errors.Internal.Wrap(err, "repository failed to get project by id")
	}
	err = p.manager.DeleteProject(project)
	if err != nil {
		return errors.Internal.Wrap(err, "project manager failed to delete project")
	}
	return p.repository.Delete(ctx, id)
}
