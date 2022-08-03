package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/mappers"
	"rol/app/validators"
	"rol/domain"
	"rol/dtos"
)

//EthernetSwitchService service structure for EthernetSwitch entity
type EthernetSwitchService struct {
	*GenericService[dtos.EthernetSwitchDto,
		dtos.EthernetSwitchCreateDto,
		dtos.EthernetSwitchUpdateDto,
		domain.EthernetSwitch]
	switchPortRepository interfaces.IGenericRepository[domain.EthernetSwitchPort]
	supportedList        *[]domain.EthernetSwitchModel
}

//NewEthernetSwitchService constructor for domain.EthernetSwitch service
//Params
//	rep - generic repository with domain.EthernetSwitch repository
//	log - logrus logger
//Return
//	New ethernet switch service
func NewEthernetSwitchService(rep interfaces.IGenericRepository[domain.EthernetSwitch], switchPortRepo interfaces.IGenericRepository[domain.EthernetSwitchPort], log *logrus.Logger) (interfaces.IGenericService[
	dtos.EthernetSwitchDto,
	dtos.EthernetSwitchCreateDto,
	dtos.EthernetSwitchUpdateDto,
	domain.EthernetSwitch], error) {
	genericService, err := NewGenericService[dtos.EthernetSwitchDto, dtos.EthernetSwitchCreateDto, dtos.EthernetSwitchUpdateDto, domain.EthernetSwitch](rep, log)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error constructing ethernet switch service")
	}
	ethernetSwitchService := &EthernetSwitchService{
		GenericService:       genericService,
		switchPortRepository: switchPortRepo,
		supportedList:        &[]domain.EthernetSwitchModel{},
	}
	ethernetSwitchService.initSupportedList()
	return ethernetSwitchService, nil
}

func (e *EthernetSwitchService) initSupportedList() {

	//Ubiquity UniFi Switch US-24-250W
	ubiquityUnifiSwitchUs24250W := domain.EthernetSwitchModel{
		Model:        "UniFi Switch US-24-250W",
		Manufacturer: "Ubiquity",
		Code:         "unifi_switch_us-24-250w",
	}
	*e.supportedList = append(*e.supportedList, ubiquityUnifiSwitchUs24250W)
}

func (e *EthernetSwitchService) sLog(ctx context.Context, level, message string) {
	entry := e.logger.WithFields(logrus.Fields{
		"actionID": ctx.Value("requestID"),
		"source":   "EthernetSwitchService",
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

func (e *EthernetSwitchService) modelIsSupported(model string) bool {
	modelIsSupported := false
	for _, supportedModel := range *e.supportedList {
		if model == supportedModel.Code {
			modelIsSupported = true
		}
	}
	return modelIsSupported
}

func (e *EthernetSwitchService) serialIsUnique(ctx context.Context, serial string, id uuid.UUID) error {
	uniqueSerialQueryBuilder := e.GenericService.repository.NewQueryBuilder(ctx)
	e.GenericService.excludeDeleted(uniqueSerialQueryBuilder)
	uniqueSerialQueryBuilder.Where("Serial", "==", serial)
	if [16]byte{} != id {
		uniqueSerialQueryBuilder.Where("ID", "!=", id)
	}
	serialEthSwitchList, err := e.GenericService.repository.GetList(ctx, "", "asc", 1, 1, uniqueSerialQueryBuilder)
	if err != nil {
		return errors.Internal.Wrap(err, "service failed get list")
	}
	if len(*serialEthSwitchList) > 0 {
		return errors.Validation.New("switch with this serial number already exist")
	}
	return nil
}

func (e *EthernetSwitchService) addressIsUnique(ctx context.Context, serial string, id uuid.UUID) error {
	uniqueSerialQueryBuilder := e.GenericService.repository.NewQueryBuilder(ctx)
	e.GenericService.excludeDeleted(uniqueSerialQueryBuilder)
	uniqueSerialQueryBuilder.Where("Address", "==", serial)
	if [16]byte{} != id {
		uniqueSerialQueryBuilder.Where("ID", "!=", id)
	}
	serialEthSwitchList, err := e.GenericService.repository.GetList(ctx, "", "asc", 1, 1, uniqueSerialQueryBuilder)
	if err != nil {
		return errors.Internal.Wrap(err, "service failed get list")
	}
	if len(*serialEthSwitchList) > 0 {
		return errors.Validation.New("switch with this address already exist")
	}
	return nil
}

//GetList Get list of ethernet switches with filtering and pagination
//Params
//	ctx - context is used only for logging
//	search - string for search in entity string fields
//	orderBy - order by entity field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return
//	*dtos.PaginatedListDto[dtos.EthernetSwitchDto] - pointer to paginated list of ethernet switches
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) GetList(ctx context.Context, search, orderBy, orderDirection string, page, pageSize int) (*dtos.PaginatedListDto[dtos.EthernetSwitchDto], error) {
	list, err := e.GenericService.GetList(ctx, search, orderBy, orderDirection, page, pageSize)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "service failed get list")
	}
	return list, nil
}

//GetByID Get ethernet switch by ID
//Params
//	ctx - context is used only for logging
//	id - entity id
//Return
//	*dtos.EthernetSwitchDto - point to ethernet switch dto
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) GetByID(ctx context.Context, id uuid.UUID) (*dtos.EthernetSwitchDto, error) {
	ethernetSwitch, err := e.GenericService.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "service failed get by id")
	}
	return ethernetSwitch, err
}

//Update save the changes to the existing ethernet switch
//Params
//	ctx - context is used only for logging
//	updateDto - ethernet switch update dto
//	id - ethernet switch id
//Return
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) Update(ctx context.Context, updateDto dtos.EthernetSwitchUpdateDto, id uuid.UUID) error {
	err := validators.ValidateEthernetSwitchUpdateDto(updateDto)
	if err != nil {
		return errors.Validation.Wrap(err, "validation failed")
	}
	if !e.modelIsSupported(updateDto.SwitchModel) {
		return errors.Validation.Wrap(err, "switch model is not supported")
	}

	err = e.serialIsUnique(ctx, updateDto.Serial, id)
	if err != nil {
		return errors.Validation.Wrap(err, "serial number uniqueness check error")
	}

	err = e.addressIsUnique(ctx, updateDto.Address, id)
	if err != nil {
		return errors.Validation.Wrap(err, "address uniqueness check error")
	}

	err = e.GenericService.Update(ctx, updateDto, id)
	if err != nil {
		return errors.Internal.Wrap(err, "service failed to update entity")
	}
	return nil
}

//Create add new ethernet switch
//Params
//	ctx - context
//	createDto - ethernet switch create dto
//Return
//	uuid.UUID - new ethernet switch id
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) Create(ctx context.Context, createDto dtos.EthernetSwitchCreateDto) (uuid.UUID, error) {
	err := validators.ValidateEthernetSwitchCreateDto(createDto)
	if err != nil {
		return [16]byte{}, errors.Validation.Wrap(err, "validation failed")
	}
	if !e.modelIsSupported(createDto.SwitchModel) {
		return [16]byte{}, errors.Validation.Wrap(err, "switch model is not supported")
	}
	err = e.serialIsUnique(ctx, createDto.Serial, [16]byte{})
	if err != nil {
		return [16]byte{}, errors.Validation.Wrap(err, "serial number uniqueness check error")
	}
	err = e.addressIsUnique(ctx, createDto.Address, [16]byte{})
	if err != nil {
		return [16]byte{}, errors.Validation.Wrap(err, "address uniqueness check error")
	}
	id, err := e.GenericService.Create(ctx, createDto)
	if err != nil {
		return [16]byte{}, errors.Internal.Wrap(err, "service failed to create entity")
	}
	return id, nil
}

//Delete mark ethernet switch as deleted
//Params
//	ctx - context is used only for logging
//	id - ethernet switch id
//Return
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) Delete(ctx context.Context, id uuid.UUID) error {
	err := e.GenericService.Delete(ctx, id)
	if err != nil {
		return errors.Internal.Wrap(err, "service failed to delete entity")
	}
	queryBuilder := e.switchPortRepository.NewQueryBuilder(ctx)
	queryBuilder.Where("EthernetSwitchID", "=", id)
	ports, err := e.switchPortRepository.GetList(ctx, "", "", 1, 100, queryBuilder)
	if err != nil {
		return errors.Internal.Wrap(err, "service failed get list")
	}
	for _, port := range *ports {
		port.SetDeleted()
		err := e.switchPortRepository.Update(ctx, &port)
		if err != nil {
			return errors.Internal.Wrap(err, "service failed to update entity")
		}
	}
	return nil
}

//GetSupportedModels Get supported switch models
//Return
//	*[]dtos.EthernetSwitchModelDto - Ethernet switch model DTO's that supported by system
func (e *EthernetSwitchService) GetSupportedModels() *[]dtos.EthernetSwitchModelDto {
	supportedModelsDtos := []dtos.EthernetSwitchModelDto{}
	for _, model := range *e.supportedList {
		modelDto := dtos.EthernetSwitchModelDto{}
		mappers.MapEthernetSwitchModelToDto(model, &modelDto)
		supportedModelsDtos = append(supportedModelsDtos, modelDto)
	}
	return &supportedModelsDtos
}
