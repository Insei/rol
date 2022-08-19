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

//EthernetSwitchPortService service structure for EthernetSwitchPort entity
type EthernetSwitchPortService struct {
	*GenericService[dtos.EthernetSwitchPortDto,
		dtos.EthernetSwitchPortCreateDto,
		dtos.EthernetSwitchPortUpdateDto,
		domain.EthernetSwitchPort]
	switchRepository interfaces.IGenericRepository[domain.EthernetSwitch]
}

//NewEthernetSwitchPortService constructor for domain.EthernetSwitchPort service
//Params
//	rep - generic repository with domain.EthernetSwitchPort repository
//	log - logrus logger
//Return
//	New ethernet switch service
func NewEthernetSwitchPortService(rep interfaces.IGenericRepository[domain.EthernetSwitchPort], switchRepo interfaces.IGenericRepository[domain.EthernetSwitch], log *logrus.Logger) (*EthernetSwitchPortService, error) {
	genericService, err := NewGenericService[dtos.EthernetSwitchPortDto, dtos.EthernetSwitchPortCreateDto, dtos.EthernetSwitchPortUpdateDto, domain.EthernetSwitchPort](rep, log)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error constructing ethernet switch port service")
	}
	ethernetSwitchPortService := &EthernetSwitchPortService{
		GenericService:   genericService,
		switchRepository: switchRepo,
	}
	return ethernetSwitchPortService, nil
}

func (e *EthernetSwitchPortService) switchIsExist(ctx context.Context, switchID uuid.UUID) (bool, error) {
	queryBuilder := e.GenericService.repository.NewQueryBuilder(ctx)
	e.GenericService.excludeDeleted(queryBuilder)
	ethernetSwitch, err := e.switchRepository.GetByIDExtended(ctx, switchID, queryBuilder)
	if err != nil {
		return false, errors.Internal.Wrap(err, "repository failed to get ethernet switch")
	}
	if ethernetSwitch == nil {
		return false, nil
	}
	return true, nil
}

func (e *EthernetSwitchPortService) sLog(ctx context.Context, level, message string) {
	entry := e.logger.WithFields(logrus.Fields{
		"actionID": ctx.Value("requestID"),
		"source":   "EthernetSwitchPortService",
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

func (e *EthernetSwitchPortService) portNameIsUniqueWithinTheSwitch(ctx context.Context, name string, switchID, id uuid.UUID) (bool, error) {
	uniqueNameQueryBuilder := e.GenericService.repository.NewQueryBuilder(ctx)
	e.GenericService.excludeDeleted(uniqueNameQueryBuilder)
	uniqueNameQueryBuilder.Where("Name", "==", name)
	uniqueNameQueryBuilder.Where("EthernetSwitchID", "==", switchID)
	if [16]byte{} != id {
		uniqueNameQueryBuilder.Where("ID", "!=", id)
	}
	ethSwitchPortsList, err := e.GenericService.repository.GetList(ctx, "", "asc", 1, 1, uniqueNameQueryBuilder)
	if err != nil {
		return false, errors.Internal.Wrap(err, "service failed to get list of switch ports")
	}
	if len(*ethSwitchPortsList) > 0 {
		return false, nil
	}
	return true, nil
}

//GetPortByID Get ethernet switch port by switch ID and port ID
//Params
//	ctx - context is used only for logging|
//  switchID - Ethernet switch ID
//	id - entity id
//Return
//	*dtos.EthernetSwitchPortDto - point to ethernet switch port dto, if existed, otherwise nil
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchPortService) GetPortByID(ctx context.Context, switchID, id uuid.UUID) (*dtos.EthernetSwitchPortDto, error) {
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error when checking the existence of the switch")
	}
	if !switchExist {
		return nil, errors.NotFound.New("switch is not found")
	}
	queryBuilder := e.repository.NewQueryBuilder(ctx)
	e.excludeDeleted(queryBuilder)
	queryBuilder.Where("EthernetSwitchID", "==", switchID)
	switchDto, err := e.getByIDBasic(ctx, id, queryBuilder)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to get by id basic")
	}
	return switchDto, nil
}

//CreatePort Create ethernet switch port by EthernetSwitchPortCreateDto
//Params
//	ctx - context is used only for logging|
//  switchID - Ethernet switch ID
//	createDto - EthernetSwitchPortCreateDto
//Return
//	*dtos.EthernetSwitchPortDto - point to ethernet switch port dto, if existed, otherwise nil
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchPortService) CreatePort(ctx context.Context, switchID uuid.UUID, createDto dtos.EthernetSwitchPortCreateDto) (uuid.UUID, error) {
	err := validators.ValidateEthernetSwitchPortCreateDto(createDto)
	if err != nil {
		return [16]byte{}, err //we already wrap error in validators
	}
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return [16]byte{}, errors.Internal.Wrap(err, "error when checking the existence of the switch")
	}
	if !switchExist {
		return [16]byte{}, errors.NotFound.New("switch is not found")
	}
	uniqName, err := e.portNameIsUniqueWithinTheSwitch(ctx, createDto.Name, switchID, [16]byte{})
	if err != nil {
		return [16]byte{}, errors.Internal.Wrap(err, "name uniqueness check error")
	}
	if !uniqName {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		return [16]byte{}, errors.AddErrorContext(err, "Name", "port with this name already exist")
	}

	entity := new(domain.EthernetSwitchPort)
	entity.EthernetSwitchID = switchID
	err = mappers.MapDtoToEntity(createDto, entity)
	if err != nil {
		return [16]byte{}, errors.Internal.Wrap(err, "failed to map ethernet switch port dto to entity")
	}

	ethernetSwitch, err := e.switchRepository.GetByIDExtended(ctx, switchID, nil)
	if err != nil {
		return [16]byte{}, errors.Internal.Wrap(err, "repository failed to get ethernet switch")
	}
	switchManager := GetEthernetSwitchManager(*ethernetSwitch)
	id, err := e.repository.Insert(ctx, *entity)
	if err != nil {
		return [16]byte{}, errors.Internal.Wrap(err, "repository failed to insert port")
	}
	if createDto.PVID > 1 {
		VLANs, _ := switchManager.GetVLANs()
		vlanExist := false
		for _, vlanID := range VLANs {
			if vlanID == createDto.PVID {
				vlanExist = true
			}
		}
		if !vlanExist {
			return [16]byte{}, errors.Internal.New("cannot set PVID for a non-existent VLAN")
		}
		err = switchManager.SetPortPVID(createDto.Name, createDto.PVID)
		if err != nil {
			return [16]byte{}, errors.Internal.Wrap(err, "set port PVID failed")
		}
	}
	return id, nil
}

//UpdatePort Update ethernet switch port
//Params
//	ctx - context is used only for logging|
//  switchID - Ethernet switch ID
//	id - entity id
//  updateDto - dtos.EthernetSwitchPortUpdateDto DTO for updating entity
//Return
//	*dtos.EthernetSwitchPortDto - point to ethernet switch port dto, if existed, otherwise nil
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchPortService) UpdatePort(ctx context.Context, switchID, id uuid.UUID, updateDto dtos.EthernetSwitchPortUpdateDto) error {
	err := validators.ValidateEthernetSwitchPortUpdateDto(updateDto)
	if err != nil {
		return err // we already wrap error in validators
	}
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return errors.Internal.Wrap(err, "error when checking the existence of the switch")
	}
	if !switchExist {
		return errors.NotFound.New("switch is not found")
	}
	uniqName, err := e.portNameIsUniqueWithinTheSwitch(ctx, updateDto.Name, switchID, id)
	if err != nil {
		return errors.Internal.Wrap(err, "name uniqueness check error")
	}
	if !uniqName {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		return errors.AddErrorContext(err, "Name", "port with this name already exist")
	}
	ethernetSwitch, err := e.switchRepository.GetByIDExtended(ctx, switchID, nil)
	if err != nil {
		return errors.Internal.Wrap(err, "repository failed to get ethernet switch")
	}
	switchManager := GetEthernetSwitchManager(*ethernetSwitch)
	queryBuilder := e.repository.NewQueryBuilder(ctx)
	e.excludeDeleted(queryBuilder)
	queryBuilder.Where("EthernetSwitchId", "==", switchID)
	err = e.updateBasic(ctx, updateDto, id, queryBuilder)
	if err != nil {
		return errors.Internal.Wrap(err, "service failed to update port")
	}
	if updateDto.PVID > 1 {
		if switchManager == nil {
			return nil
		}
		VLANs, _ := switchManager.GetVLANs()
		vlanExist := false
		for _, vlanID := range VLANs {
			if vlanID == updateDto.PVID {
				vlanExist = true
			}
		}
		if !vlanExist {
			return errors.Internal.New("cannot set PVID for a non-existent VLAN")
		}
		err = switchManager.SetPortPVID(updateDto.Name, updateDto.PVID)
		if err != nil {
			return errors.Internal.Wrap(err, "set port PVID failed")
		}
	}
	return nil
}

//GetPorts Get list of ethernet switch ports with filtering and pagination
//Params
//	ctx - context is used only for logging
//  switchID - uuid of the ethernet switch
//	search - string for search in ethernet switch port string fields
//	orderBy - order by ethernet switch port field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return
//	*dtos.PaginatedListDto[dtos.EthernetSwitchPortDto] - pointer to paginated list of ethernet switches
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchPortService) GetPorts(ctx context.Context, switchID uuid.UUID, search, orderBy, orderDirection string, page, pageSize int) (*dtos.PaginatedListDto[dtos.EthernetSwitchPortDto], error) {
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error when checking the existence of the switch")
	}
	if !switchExist {
		return nil, errors.NotFound.New("switch is not found")
	}
	queryBuilder := e.repository.NewQueryBuilder(ctx)
	e.excludeDeleted(queryBuilder)
	queryBuilder.Where("EthernetSwitchId", "==", switchID)
	if len(search) > 3 {
		e.addSearchInAllFields(search, queryBuilder)
	}
	return e.getListBasic(ctx, queryBuilder, orderBy, orderDirection, page, pageSize)
}

//DeletePort mark ethernet switch port as deleted
//Params
//	ctx - context is used only for logging
//	switchID - ethernet switch id
//Return
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchPortService) DeletePort(ctx context.Context, switchID, id uuid.UUID) error {
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return errors.Internal.Wrap(err, "error when checking the existence of the switch")
	}
	if !switchExist {
		return errors.NotFound.New("switch is not found")
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

//func (e *EthernetSwitchPortService) createVLANsOnPort(ethernetSwitch domain.EthernetSwitch, dto dtos.EthernetSwitchPortBaseDto) error {
//	switchManager := GetEthernetSwitchManager(ethernetSwitch)
//	if switchManager == nil {
//		return nil
//	}
//	for _, taggedVLANID := range dto.TaggedVLANs {
//		VLANFound, err := e.isVLANExistOnSwitch(taggedVLANID, switchManager)
//		if err != nil {
//			return errors.Internal.Wrap(err, "error check VLAN existence")
//		}
//		if !VLANFound {
//			err := switchManager.CreateVLAN(taggedVLANID)
//			if err != nil {
//				return errors.Internal.Wrap(err, "error create tagged VLAN")
//			}
//		}
//		err = switchManager.AddTaggedVLANOnPort(dto.Name, taggedVLANID)
//		if err != nil {
//			return errors.Internal.Wrap(err, "error add tagged VLAN")
//		}
//	}
//	VLANFound, err := e.isVLANExistOnSwitch(dto.UntaggedVLAN, switchManager)
//	if err != nil {
//		return errors.Internal.Wrap(err, "error check VLAN existence")
//	}
//	if !VLANFound {
//		err := switchManager.CreateVLAN(dto.UntaggedVLAN)
//		if err != nil {
//			return errors.Internal.Wrap(err, "error create tagged VLAN")
//		}
//	}
//	if dto.UntaggedVLAN > 0 {
//		err := switchManager.SetUntaggedVLANOnPort(dto.Name, dto.UntaggedVLAN)
//		if err != nil {
//			return errors.Internal.Wrap(err, "error create untagged VLAN")
//		}
//	}
//	return nil
//}
//
//func (e *EthernetSwitchPortService) updateVLANsOnPort(ethernetSwitch domain.EthernetSwitch, dto dtos.EthernetSwitchPortBaseDto) error {
//	switchManager := GetEthernetSwitchManager(ethernetSwitch)
//	if switchManager == nil {
//		return nil
//	}
//	_, tVLANs, err := switchManager.GetVLANsOnPort(dto.Name)
//	if err != nil {
//		return errors.Internal.Wrap(err, "error get VLANs")
//	}
//
//	for _, VLANOnPort := range tVLANs {
//		err = switchManager.RemoveTaggedVLANOnPort(dto.Name, VLANOnPort)
//		if err != nil {
//			return errors.Internal.Wrap(err, "error remove VLAN")
//		}
//	}
//	for _, taggedVLAN := range dto.TaggedVLANs {
//		err = switchManager.AddTaggedVLANOnPort(dto.Name, taggedVLAN)
//		if err != nil {
//			return errors.Internal.Wrap(err, "error add VLAN")
//		}
//	}
//	err = switchManager.AddTaggedVLANOnPort(dto.Name, dto.UntaggedVLAN)
//	if err != nil {
//		return errors.Internal.Wrap(err, "error add VLAN")
//	}
//	return nil
//}
//
//func (e *EthernetSwitchPortService) isVLANExistOnSwitch(id int, switchManager interfaces.IEthernetSwitchManager) (bool, error) {
//	VLANs, err := switchManager.GetVLANs()
//	if err != nil {
//		return false, errors.Internal.Wrap(err, "error get VLANs")
//	}
//	exists := false
//	for _, VLANOnSwitch := range VLANs {
//		if id == VLANOnSwitch {
//			exists = true
//		}
//	}
//	return exists, nil
//}
