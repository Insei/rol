package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"reflect"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/mappers"
	"rol/domain"
	"rol/dtos"
)

//TFTPServerService service structure for TFTP server
type TFTPServerService struct {
	configsRepo interfaces.IGenericRepository[uuid.UUID, domain.TFTPConfig]
	pathsRepo   interfaces.IGenericRepository[uuid.UUID, domain.TFTPPathRatio]
	factory     interfaces.ITFTPServerFactory
	servers     map[uuid.UUID]interfaces.ITFTPServer
	logger      *logrus.Logger
	//logSourceName - logger recording source
	logSourceName string
}

//NewTFTPServerService constructor for TFTP server service
//
//Params
//	configsRepo - generic repository with domain.TFTPConfig entity
//	pathsRepo - generic repository with domain.TFTPPathRatio entity
//	manager - tftp server manager
//Return
//	New TFTP server service
func NewTFTPServerService(configsRepo interfaces.IGenericRepository[uuid.UUID, domain.TFTPConfig],
	pathsRepo interfaces.IGenericRepository[uuid.UUID, domain.TFTPPathRatio],
	factory interfaces.ITFTPServerFactory, log *logrus.Logger) *TFTPServerService {
	return &TFTPServerService{
		configsRepo:   configsRepo,
		pathsRepo:     pathsRepo,
		factory:       factory,
		logger:        log,
		logSourceName: reflect.TypeOf(TFTPServerService{}).Name(),
		servers:       map[uuid.UUID]interfaces.ITFTPServer{},
	}
}

func (s *TFTPServerService) log(ctx context.Context, level, message string) {
	if ctx != nil {
		actionID := uuid.UUID{}
		if ctx.Value("requestID") != nil {
			actionID = ctx.Value("requestID").(uuid.UUID)
		}

		entry := s.logger.WithFields(logrus.Fields{
			"actionID": actionID,
			"source":   s.logSourceName,
		})
		switch level {
		case "err", "error":
			entry.Error(message)
		case "info":
			entry.Info(message)
		case "warn", "warning":
			entry.Warn(message)
		case "debug":
			entry.Debug(message)
		}
	}
}

func (s *TFTPServerService) getServerState(configID uuid.UUID) domain.TFTPServerState {
	if server, ok := s.servers[configID]; ok {
		return server.GetState()
	}
	return domain.TFTPStateStopped
}

func (s *TFTPServerService) getServerPaths(ctx context.Context, id uuid.UUID) ([]domain.TFTPPathRatio, error) {
	paths := []domain.TFTPPathRatio{}
	queryBuilder := s.getQueryBuilderWithConfigID(ctx, id)
	pathsCount, err := s.pathsRepo.Count(ctx, queryBuilder)
	if err != nil {
		return paths, errors.Internal.Wrap(err, "failed to count tftp paths")
	}
	return s.pathsRepo.GetList(ctx, "", "", 1, int(pathsCount), queryBuilder)
}

//TFTPServerServiceInit initialize TFTP service
func TFTPServerServiceInit(s *TFTPServerService) error {
	ctx := context.Background()
	queryBuilder := s.configsRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("Enabled", "==", true)
	enabledServersCount, err := s.configsRepo.Count(ctx, queryBuilder)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to count enabled tftp servers configs")
	}
	enabledServersConfigs, err := s.configsRepo.GetList(ctx, "", "", 1, int(enabledServersCount), nil)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to get tftp server configs list")
	}
	for _, config := range enabledServersConfigs {
		server, err := s.factory.Create(config)
		if err != nil {
			return errors.Internal.Wrapf(err, "failed to create runtime tftp server with id %s", config.ID)
		}
		s.servers[config.ID] = server
		err = server.Start()
		if err != nil {
			return errors.Internal.Wrapf(err, "failed to start runtime tftp server with id %s", config.ID)
		}
	}
	return nil
}

//GetServerByID get TFTP server by ID
//
//Params
//	ctx - context
//	id - TFTP server id
//Return
//	dtos.TFTPConfigDto - TFTP server dto
//	error - if an error occurs, otherwise nil
func (s *TFTPServerService) GetServerByID(ctx context.Context, id uuid.UUID) (dtos.TFTPServerDto, error) {
	dto, err := GetByID[dtos.TFTPServerDto](ctx, s.configsRepo, id, nil)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "failed to get tftp server config")
	}
	dto.State = s.getServerState(dto.ID).String()
	return dto, err
}

//GetServerList get list of TFTP servers with filtering and pagination
//
//Params
//	ctx - context
//	search - string for search in entity string fields
//	orderBy - order by entity field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return
//	dtos.PaginatedItemsDto[dtos.TFTPConfigDto] - paginated list of TFTP servers
//	error - if an error occurs, otherwise nil
func (s *TFTPServerService) GetServerList(ctx context.Context, search, orderBy, orderDirection string, page, pageSize int) (dtos.PaginatedItemsDto[dtos.TFTPServerDto], error) {
	servers, err := GetList[dtos.TFTPServerDto](ctx, s.configsRepo, search, orderBy, orderDirection, page, pageSize)
	if err != nil {
		return servers, errors.Internal.Wrap(err, "failed to get tftp server configs list")
	}
	for i, server := range servers.Items {
		(&servers.Items[i]).State = s.getServerState(server.ID).String()
	}
	return servers, nil
}

//CreateServer add new TFTP server
//
//Params
//	ctx - context
//	createDto - TFTP server create dto
//Return
//	dtos.TFTPConfigDto - created TFTP server
//	error - if an error occurs, otherwise nil
func (s *TFTPServerService) CreateServer(ctx context.Context, createDto dtos.TFTPServerCreateDto) (dtos.TFTPServerDto, error) {
	//prepare dto and entity
	dto := dtos.TFTPServerDto{}
	entity := new(domain.TFTPConfig)
	//map create config dto fields to config entity
	err := mappers.MapDtoToEntity(createDto, entity)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error map entity to dto")
	}
	//write tftp config entity to repository
	config, err := s.configsRepo.Insert(ctx, *entity)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "create tftp server config error")
	}
	//start runtime tftp server and set state to dto.State
	dto.State = domain.TFTPStateStopped.String()
	if config.Enabled {
		//create tftp server and if all is ok,
		//start this server and place server to servers pool
		//otherwise ignore
		server, err := s.factory.Create(config)
		if err == nil {
			err = server.Start()
			if err == nil {
				s.servers[config.ID] = server
				dto.State = server.GetState().String()
			}
		}
	}
	//map config fields to dto
	err = mappers.MapEntityToDto(config, &dto)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error map dto to entity")
	}
	return dto, nil
}

//UpdateServer save the changes to the existing TFTP server
//
//Params
//	ctx - context is used only for logging
//	updateDto - TFTP server update dto
//	id - TFTP server id
//Return
//	dtos.TFTPConfigDto - updated TFTP server
//	error - if an error occurs, otherwise nil
func (s *TFTPServerService) UpdateServer(ctx context.Context, updateDto dtos.TFTPServerUpdateDto, id uuid.UUID) (dtos.TFTPServerDto, error) {
	dto := dtos.TFTPServerDto{}
	config, err := s.configsRepo.GetByID(ctx, id)
	if err != nil {
		return dto, err
	}
	err = mappers.MapDtoToEntity(updateDto, &config)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error map dto to entity")
	}
	updatedConfig, err := s.configsRepo.Update(ctx, config)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "failed to update entity in repository")
	}
	//tftp runtime servers management
	dto.State = domain.TFTPStateStopped.String()
	if updatedConfig.Enabled {
		//create runtime tftp server if not exist in server pool
		if server, ok := s.servers[updatedConfig.ID]; !ok {
			server, err = s.factory.Create(config)
			if err != nil {
				s.servers[updatedConfig.ID] = server
			}
		}
		//start server if this server exist in server pool
		if server, ok := s.servers[updatedConfig.ID]; ok {
			_ = server.Start()
			dto.State = server.GetState().String()
		}
	} else {
		if server, ok := s.servers[updatedConfig.ID]; ok {
			server.Stop()
			delete(s.servers, updatedConfig.ID)
		}
	}
	err = mappers.MapEntityToDto(updatedConfig, &dto)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error map tftp config to server dto")
	}
	return dto, nil
}

//DeleteServer delete TFTP server
//
//Params
//	ctx - context is used only for logging
//	id - TFTP server id
//Return
//	error - if an error occurs, otherwise nil
func (s *TFTPServerService) DeleteServer(ctx context.Context, id uuid.UUID) error {
	exist, err := s.configsRepo.IsExist(ctx, id, nil)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to check existence of tftp server")
	}
	if !exist {
		return errors.NotFound.New("tftp server with this id is not found")
	}
	err = s.pathsRepo.DeleteAll(ctx, s.getQueryBuilderWithConfigID(ctx, id))
	if err != nil {
		return errors.Internal.Wrap(err, "failed to remove all tftp server paths ratios")
	}
	err = s.configsRepo.Delete(ctx, id)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to delete tftp server config")
	}
	if server, ok := s.servers[id]; ok {
		server.Stop()
		delete(s.servers, id)
	}
	return nil
}

func (s *TFTPServerService) updateRuntimeServerPaths(ctx context.Context, id uuid.UUID) error {
	var (
		server interfaces.ITFTPServer
		ok     bool
	)
	if server, ok = s.servers[id]; !ok {
		return nil
	}
	paths, err := s.getServerPaths(ctx, id)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to get TFTP server paths")
	}
	err = server.ReloadPaths(paths)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to reload TFTP server paths")
	}
	return nil
}

func (s *TFTPServerService) getQueryBuilderWithConfigID(ctx context.Context, tftpConfigID uuid.UUID) interfaces.IQueryBuilder {
	queryBuilder := s.pathsRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("TFTPConfigID", "==", tftpConfigID)
	return queryBuilder
}

//GetPathsList get list of TFTP servers paths with filtering and pagination
//
//Params
//	ctx - context is used only for logging
//	configID - tftp config id
//	orderBy - order by entity field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return
//	dtos.PaginatedItemsDto[dtos.TFTPPathDto] - paginated list of TFTP servers paths
//	error - if an error occurs, otherwise nil
func (s *TFTPServerService) GetPathsList(ctx context.Context, configID uuid.UUID, orderBy, orderDirection string, page, pageSize int) (dtos.PaginatedItemsDto[dtos.TFTPPathDto], error) {
	return GetListExtended[dtos.TFTPPathDto](ctx, s.pathsRepo, s.getQueryBuilderWithConfigID(ctx, configID), orderBy, orderDirection, page, pageSize)
}

//GetPathByID get list of TFTP servers paths with filtering and pagination
//
//Params
//	ctx - context
//	configID - TFTP server id
//	pathID - TFTP path id
//Return
//	dtos.TFTPPathDto - TFTP server dto
//	error - if an error occurs, otherwise nil
func (s *TFTPServerService) GetPathByID(ctx context.Context, configID uuid.UUID, pathID uuid.UUID) (dtos.TFTPPathDto, error) {
	return GetByID[dtos.TFTPPathDto](ctx, s.pathsRepo, pathID, s.getQueryBuilderWithConfigID(ctx, configID))
}

//CreatePath add new TFTP server path
//
//Params
//	ctx - context
//	configID - tftp config id
//	createDto - TFTP server create dto
//Return
//	dtos.TFTPPathDto - created TFTP path
//	error - if an error occurs, otherwise nil
func (s *TFTPServerService) CreatePath(ctx context.Context, configID uuid.UUID, createDto dtos.TFTPPathCreateDto) (dtos.TFTPPathDto, error) {
	entity := new(domain.TFTPPathRatio)
	outDto := new(dtos.TFTPPathDto)
	err := mappers.MapDtoToEntity(createDto, entity)
	if err != nil {
		return *outDto, errors.Internal.Wrap(err, "error map entity to dto")
	}
	entity.TFTPConfigID = configID
	newEntity, err := s.pathsRepo.Insert(ctx, *entity)
	if err != nil {
		return *outDto, errors.Internal.Wrap(err, "create entity error")
	}
	err = s.updateRuntimeServerPaths(ctx, configID)
	if err != nil {
		s.log(ctx, "err", err.Error())
	}
	err = mappers.MapEntityToDto(newEntity, outDto)
	if err != nil {
		return *outDto, errors.Internal.Wrap(err, "error map dto to entity")
	}
	return *outDto, nil
}

func (s *TFTPServerService) existencePathCheck(ctx context.Context, configID uuid.UUID, id uuid.UUID) error {
	exist, err := s.pathsRepo.IsExist(ctx, id, s.getQueryBuilderWithConfigID(ctx, configID))
	if err != nil {
		return errors.Internal.Wrapf(err,
			"failed to get path ratio for tftp config id: %s, path id: %s", configID.String(), id.String())
	}
	if !exist {
		return errors.NotFound.Wrap(err, "path ratio not found on tftp server")
	}
	return nil
}

//UpdatePath update TFTP server path
//
//Params
//	ctx - context
//	configID - tftp config id
//	pathID - tftp path id
//	updateDto - TFTP path update dto
//Return
//	dtos.TFTPPathDto - updated TFTP path
//	error - if an error occurs, otherwise nil
func (s *TFTPServerService) UpdatePath(ctx context.Context, configID, pathID uuid.UUID, updateDto dtos.TFTPPathUpdateDto) (dtos.TFTPPathDto, error) {
	dto := new(dtos.TFTPPathDto)
	pathRatio, err := s.pathsRepo.GetByIDExtended(ctx, pathID, s.getQueryBuilderWithConfigID(ctx, configID))
	if err != nil {
		return *dto, err
	}
	err = mappers.MapDtoToEntity(updateDto, pathRatio)
	if err != nil {
		return *dto, errors.Internal.Wrap(err, "error map entity to dto")
	}
	updatedPathRatio, err := s.pathsRepo.Update(ctx, pathRatio)
	if err != nil {
		return *dto, errors.Internal.Wrap(err, "update tftp path ratio error")
	}
	err = s.updateRuntimeServerPaths(ctx, configID)
	if err != nil {
		s.log(ctx, "error", err.Error())
	}
	err = mappers.MapEntityToDto(updatedPathRatio, dto)
	if err != nil {
		return *dto, errors.Internal.Wrap(err, "error map dto to entity")
	}
	return *dto, nil
}

//DeletePath mark TFTP server path as deleted
//
//Params
//	ctx - context is used only for logging
//	configID - id of tftp config
//	id - TFTP server path id
//Return
//	error - if an error occurs, otherwise nil
func (s *TFTPServerService) DeletePath(ctx context.Context, configID uuid.UUID, id uuid.UUID) error {
	err := s.existencePathCheck(ctx, configID, id)
	if err != nil {
		return err
	}
	err = s.pathsRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	err = s.updateRuntimeServerPaths(ctx, configID)
	if err != nil {
		s.log(ctx, "err", err.Error())
	}
	return nil
}
