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

const (
	//ErrorGetSwitch repository failed to get ethernet switch
	ErrorGetSwitch = "repository failed to get ethernet switch"
	//ErrorSwitchExistence error when checking the existence of the switch
	ErrorSwitchExistence = "error when checking the existence of the switch"
	//ErrorSwitchNotFound switch is not found error
	ErrorSwitchNotFound = "switch is not found"
	//ErrorGetPortByID get port by id failed
	ErrorGetPortByID = "get port by id failed"
	//ErrorAddTaggedVLAN add tagged VLAN on port failed
	ErrorAddTaggedVLAN = "add tagged VLAN on port failed"
	//ErrorRemoveVLAN failed to remove VLAN from port
	ErrorRemoveVLAN = "failed to remove VLAN from port"
)

//EthernetSwitchVLANService service structure for EthernetSwitchPort entity
type EthernetSwitchVLANService struct {
	*GenericService[dtos.EthernetSwitchVLANDto,
		dtos.EthernetSwitchVLANCreateDto,
		dtos.EthernetSwitchVLANUpdateDto,
		domain.EthernetSwitchVLAN]
	switchRepository     interfaces.IGenericRepository[domain.EthernetSwitch]
	switchPortRepository interfaces.IGenericRepository[domain.EthernetSwitchPort]
}

//NewEthernetSwitchVLANService constructor for domain.EthernetSwitchPort service
//Params
//	rep - generic repository with domain.EthernetSwitchPort repository
//	log - logrus logger
//Return
//	New ethernet switch service
func NewEthernetSwitchVLANService(rep interfaces.IGenericRepository[domain.EthernetSwitchVLAN], switchRepo interfaces.IGenericRepository[domain.EthernetSwitch], switchPortRepo interfaces.IGenericRepository[domain.EthernetSwitchPort], log *logrus.Logger) (*EthernetSwitchVLANService, error) {
	genericService, err := NewGenericService[dtos.EthernetSwitchVLANDto, dtos.EthernetSwitchVLANCreateDto, dtos.EthernetSwitchVLANUpdateDto, domain.EthernetSwitchVLAN](rep, log)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error constructing ethernet switch port service")
	}
	ethernetSwitchVLANService := &EthernetSwitchVLANService{
		GenericService:       genericService,
		switchRepository:     switchRepo,
		switchPortRepository: switchPortRepo,
	}
	return ethernetSwitchVLANService, nil
}

func (e *EthernetSwitchVLANService) switchIsExist(ctx context.Context, switchID uuid.UUID) (bool, error) {
	queryBuilder := e.GenericService.repository.NewQueryBuilder(ctx)
	e.GenericService.excludeDeleted(queryBuilder)
	ethernetSwitch, err := e.switchRepository.GetByIDExtended(ctx, switchID, queryBuilder)
	if err != nil {
		return false, errors.Internal.Wrap(err, ErrorGetSwitch)
	}
	if ethernetSwitch == nil {
		return false, nil
	}
	return true, nil
}

func (e *EthernetSwitchVLANService) isVLANIdUnique(ctx context.Context, vlanID int, switchID uuid.UUID) (bool, error) {
	uniqueIDQueryBuilder := e.GenericService.repository.NewQueryBuilder(ctx)
	e.GenericService.excludeDeleted(uniqueIDQueryBuilder)
	uniqueIDQueryBuilder.Where("VlanID", "==", vlanID)
	uniqueIDQueryBuilder.Where("EthernetSwitchID", "==", switchID)
	switchVLANsList, err := e.GenericService.repository.GetList(ctx, "", "asc", 1, 1, uniqueIDQueryBuilder)
	if err != nil {
		return false, errors.Internal.Wrap(err, "service failed to get list of switch ports")
	}
	if len(*switchVLANsList) > 0 {
		return false, nil
	}

	return true, nil
}

//GetVLANByID Get ethernet switch VLAN by switch ID and VLAN ID
//
//Params
//	ctx - context is used only for logging|
//  switchID - Ethernet switch ID
//	id - VLAN id
//Return
//	*dtos.EthernetSwitchVLANDto - point to ethernet switch VLAN dto, if existed, otherwise nil
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchVLANService) GetVLANByID(ctx context.Context, switchID, id uuid.UUID) (*dtos.EthernetSwitchVLANDto, error) {
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return nil, errors.Internal.Wrap(err, ErrorSwitchExistence)
	}
	if !switchExist {
		return nil, errors.NotFound.New(ErrorSwitchNotFound)
	}
	queryBuilder := e.repository.NewQueryBuilder(ctx)
	e.excludeDeleted(queryBuilder)
	queryBuilder.Where("EthernetSwitchID", "==", switchID)
	VLANDto, err := e.getByIDBasic(ctx, id, queryBuilder)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to get by id basic")
	}
	return VLANDto, nil
}

//GetVLANs Get list of ethernet switch VLANs with filtering and pagination
//
//Params
//	ctx - context is used only for logging
//  switchID - uuid of the ethernet switch
//	search - string for search in ethernet switch port string fields
//	orderBy - order by ethernet switch port field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return
//	*dtos.PaginatedItemsDto[dtos.EthernetSwitchVLANDto] - pointer to paginated list of ethernet switch VLANs
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchVLANService) GetVLANs(ctx context.Context, switchID uuid.UUID, search, orderBy, orderDirection string, page, pageSize int) (*dtos.PaginatedItemsDto[dtos.EthernetSwitchVLANDto], error) {
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return nil, errors.Internal.Wrap(err, ErrorSwitchExistence)
	}
	if !switchExist {
		return nil, errors.NotFound.New(ErrorSwitchNotFound)
	}
	queryBuilder := e.repository.NewQueryBuilder(ctx)
	e.excludeDeleted(queryBuilder)
	queryBuilder.Where("EthernetSwitchID", "==", switchID)
	if len(search) > 3 {
		e.addSearchInAllFields(search, queryBuilder)
	}
	return e.getListBasic(ctx, queryBuilder, orderBy, orderDirection, page, pageSize)
}

//CreateVLAN Create ethernet switch VLAN by EthernetSwitchVLANCreateDto
//Params
//	ctx - context is used only for logging
//  switchID - Ethernet switch ID
//	createDto - EthernetSwitchVLANCreateDto
//Return
//	dtos.EthernetSwitchVLANDto - created switch VLAN
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchVLANService) CreateVLAN(ctx context.Context, switchID uuid.UUID, createDto dtos.EthernetSwitchVLANCreateDto) (dtos.EthernetSwitchVLANDto, error) {
	err := validators.ValidateEthernetSwitchVLANCreateDto(createDto)
	if err != nil {
		return dtos.EthernetSwitchVLANDto{}, err //we already wrap error in validators
	}
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, ErrorSwitchExistence)
	}
	if !switchExist {
		return dtos.EthernetSwitchVLANDto{}, errors.NotFound.New(ErrorSwitchNotFound)
	}
	uniqVLANId, err := e.isVLANIdUnique(ctx, createDto.VlanID, switchID)
	if err != nil {
		return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, "VLAN ID uniqueness check error")
	}
	if !uniqVLANId {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		return dtos.EthernetSwitchVLANDto{}, errors.AddErrorContext(err, "VlanID", "vlan with this id already exist")
	}
	ethernetSwitch, err := e.switchRepository.GetByIDExtended(ctx, switchID, nil)
	if err != nil {
		return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, ErrorGetSwitch)
	}
	entity := new(domain.EthernetSwitchVLAN)
	err = mappers.MapDtoToEntity(createDto, entity)
	entity.EthernetSwitchID = switchID
	if err != nil {
		return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, "failed to map ethernet switch port dto to entity")
	}
	newEntity, err := e.repository.Insert(ctx, *entity)
	if err != nil {
		return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, "repository failed to insert VLAN")
	}
	dto := dtos.EthernetSwitchVLANDto{}
	switchManager := GetEthernetSwitchManager(*ethernetSwitch)
	if switchManager == nil {
		err = mappers.MapEntityToDto(newEntity, &dto)
		if err != nil {
			return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, "error map entity to dto")
		}
		return dto, nil
	}
	err = switchManager.CreateVLAN(createDto.VlanID)
	if err != nil {
		return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, "create VLAN on switch failed")
	}
	for _, taggedPortID := range createDto.TaggedPorts {
		switchPort, err := e.switchPortRepository.GetByID(ctx, taggedPortID)
		if err != nil {
			return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, ErrorGetPortByID)
		}
		err = switchManager.AddTaggedVLANOnPort(switchPort.Name, createDto.VlanID)
		if err != nil {
			return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, ErrorAddTaggedVLAN)
		}
	}
	for _, untaggedPortID := range createDto.UntaggedPorts {
		switchPort, err := e.switchPortRepository.GetByID(ctx, untaggedPortID)
		if err != nil {
			return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, ErrorGetPortByID)
		}
		err = switchManager.AddUntaggedVLANOnPort(switchPort.Name, createDto.VlanID)
		if err != nil {
			return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, ErrorAddTaggedVLAN)
		}
	}
	err = switchManager.SaveConfig()
	if err != nil {
		return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, "save switch config failed")
	}
	err = mappers.MapEntityToDto(newEntity, &dto)
	if err != nil {
		return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, "error map entity to dto")
	}
	return dto, nil
}

//UpdateVLAN Update ethernet switch VLAN
//
//Params
//	ctx - context is used only for logging
//  switchID - Ethernet switch ID
//	id - VLAN id
//  updateDto - dtos.EthernetSwitchVLANUpdateDto DTO for updating entity
//Return
//	dtos.EthernetSwitchVLANDto - updated switch VLAN
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchVLANService) UpdateVLAN(ctx context.Context, switchID, id uuid.UUID, updateDto dtos.EthernetSwitchVLANUpdateDto) (dtos.EthernetSwitchVLANDto, error) {
	err := validators.ValidateEthernetSwitchVLANUpdateDto(updateDto)
	if err != nil {
		return dtos.EthernetSwitchVLANDto{}, err //we already wrap error in validators
	}
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, ErrorSwitchExistence)
	}
	if !switchExist {
		return dtos.EthernetSwitchVLANDto{}, errors.NotFound.New(ErrorSwitchNotFound)
	}

	VLAN, err := e.GetVLANByID(ctx, switchID, id)
	if err != nil {
		return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, "get VLAN by id failed")
	}
	if VLAN == nil {
		return dtos.EthernetSwitchVLANDto{}, errors.NotFound.New("VLAN not found")
	}
	ethernetSwitch, err := e.switchRepository.GetByIDExtended(ctx, switchID, nil)
	if err != nil {
		return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, ErrorGetSwitch)
	}
	switchManager := GetEthernetSwitchManager(*ethernetSwitch)
	if switchManager != nil {
		//TODO: warning message log that switch manager is nil
		err = e.updateVLANsOnPort(ctx, *VLAN, updateDto, switchManager)
		if err != nil {
			return dtos.EthernetSwitchVLANDto{}, errors.Internal.Wrap(err, "update VLANs on port failed")
		}
	}

	queryBuilder := e.repository.NewQueryBuilder(ctx)
	e.excludeDeleted(queryBuilder)
	queryBuilder.Where("EthernetSwitchID", "==", switchID)
	return e.updateBasic(ctx, updateDto, id, queryBuilder)
}

//DeleteVLAN mark ethernet switch VLAN as deleted
//Params
//	ctx - context is used only for logging
//	switchID - ethernet switch id
//	id - VLAN ID
//Return
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchVLANService) DeleteVLAN(ctx context.Context, switchID, id uuid.UUID) error {
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorSwitchExistence)
	}
	if !switchExist {
		return errors.NotFound.New(ErrorSwitchNotFound)
	}
	queryBuilder := e.repository.NewQueryBuilder(ctx)
	e.excludeDeleted(queryBuilder)
	queryBuilder.Where("EthernetSwitchID", "==", switchID)
	entity, err := e.repository.GetByIDExtended(ctx, id, queryBuilder)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to get by id")
	}
	if entity == nil {
		return errors.NotFound.New("ethernet switch port is not exist")
	}
	return e.GenericService.Delete(ctx, id)
}

func (e *EthernetSwitchVLANService) updateVLANsOnPort(ctx context.Context, VLAN dtos.EthernetSwitchVLANDto, updateDto dtos.EthernetSwitchVLANUpdateDto, switchManager interfaces.IEthernetSwitchManager) error {
	diffTaggedToRemove := e.getDifference(VLAN.TaggedPorts, updateDto.TaggedPorts)
	for _, id := range diffTaggedToRemove {
		switchPort, err := e.switchPortRepository.GetByID(ctx, id)
		if err != nil {
			return errors.Internal.Wrap(err, ErrorGetPortByID)
		}
		err = switchManager.RemoveVLANFromPort(switchPort.Name, VLAN.VlanID)
		if err != nil {
			return errors.Internal.Wrap(err, ErrorRemoveVLAN)
		}
	}
	diffTaggedToAdd := e.getDifference(updateDto.TaggedPorts, VLAN.TaggedPorts)
	for _, id := range diffTaggedToAdd {
		switchPort, err := e.switchPortRepository.GetByID(ctx, id)
		if err != nil {
			return errors.Internal.Wrap(err, ErrorGetPortByID)
		}
		err = switchManager.AddTaggedVLANOnPort(switchPort.Name, VLAN.VlanID)
		if err != nil {
			return errors.Internal.Wrap(err, "failed to add tagged VLAN on port")
		}
	}
	diffUntaggedToRemove := e.getDifference(VLAN.UntaggedPorts, updateDto.UntaggedPorts)
	for _, id := range diffUntaggedToRemove {
		switchPort, err := e.switchPortRepository.GetByID(ctx, id)
		if err != nil {
			return errors.Internal.Wrap(err, ErrorGetPortByID)
		}
		err = switchManager.RemoveVLANFromPort(switchPort.Name, VLAN.VlanID)
		if err != nil {
			return errors.Internal.Wrap(err, ErrorRemoveVLAN)
		}
	}
	diffUntaggedToAdd := e.getDifference(VLAN.UntaggedPorts, updateDto.UntaggedPorts)
	for _, id := range diffUntaggedToAdd {
		switchPort, err := e.switchPortRepository.GetByID(ctx, id)
		if err != nil {
			return errors.Internal.Wrap(err, ErrorGetPortByID)
		}
		err = switchManager.AddUntaggedVLANOnPort(switchPort.Name, VLAN.VlanID)
		if err != nil {
			return errors.Internal.Wrap(err, "failed to add untagged VLAN on port")
		}
	}
	return nil
}

func (e *EthernetSwitchVLANService) getDifference(a, b []uuid.UUID) []uuid.UUID {
	mb := make(map[uuid.UUID]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []uuid.UUID
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
