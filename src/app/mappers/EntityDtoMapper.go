package mappers

import (
	"rol/app/errors"
	"rol/domain"
	"rol/dtos"
)

//MapDtoToEntity map a DTO to its corresponding entity
//Params
//	dto - DTO struct
//	entity - pointer to the entity
//Return
//  error - if error occurs return error, otherwise nil
func MapDtoToEntity(dto interface{}, entity interface{}) error {
	switch dto.(type) {
	// EthernetSwitch
	case dtos.EthernetSwitchCreateDto:
		MapEthernetSwitchCreateDto(dto.(dtos.EthernetSwitchCreateDto), entity.(*domain.EthernetSwitch))
	case dtos.EthernetSwitchUpdateDto:
		MapEthernetSwitchUpdateDto(dto.(dtos.EthernetSwitchUpdateDto), entity.(*domain.EthernetSwitch))
	//EthernetSwitchPort
	case dtos.EthernetSwitchPortCreateDto:
		MapEthernetSwitchPortCreateDto(dto.(dtos.EthernetSwitchPortCreateDto), entity.(*domain.EthernetSwitchPort))
	case dtos.EthernetSwitchPortUpdateDto:
		MapEthernetSwitchPortUpdateDto(dto.(dtos.EthernetSwitchPortUpdateDto), entity.(*domain.EthernetSwitchPort))
	//HostNetworkVlan
	case dtos.HostNetworkVlanCreateDto:
		MapHostNetworkCreateDtoToEntity(dto.(dtos.HostNetworkVlanCreateDto), entity.(*domain.HostNetworkVlan))
	case dtos.HostNetworkVlanUpdateDto:
		MapHostNetworkUpdateDtoToEntity(dto.(dtos.HostNetworkVlanUpdateDto), entity.(*domain.HostNetworkVlan))
	//TFTPConfig
	case dtos.TFTPConfigCreateDto:
		MapTFTPConfigCreateDtoToEntity(dto.(dtos.TFTPConfigCreateDto), entity.(*domain.TFTPConfig))
	case dtos.TFTPConfigUpdateDto:
		MapTFTPConfigUpdateDtoToEntity(dto.(dtos.TFTPConfigUpdateDto), entity.(*domain.TFTPConfig))
	//TFTPPathRatio
	case dtos.TFTPPathCreateDto:
		MapTFTPPathCreateDtoToEntity(dto.(dtos.TFTPPathCreateDto), entity.(*domain.TFTPPathRatio))
	default:
		return errors.Internal.Newf("can't find route for map dto %+v to entity %+v", dto, entity)
	}
	return nil
}

//MapEntityToDto map entity to DTO
//Params
//	entity - Entity struct
//	dto - dest DTO
//Return
//  error - if error occurs return error, otherwise nil
func MapEntityToDto(entity interface{}, dto interface{}) error {
	switch entity.(type) {
	// EthernetSwitch
	case domain.EthernetSwitch:
		MapEthernetSwitchToDto(entity.(domain.EthernetSwitch), dto.(*dtos.EthernetSwitchDto))
	//EthernetSwitchPort
	case domain.EthernetSwitchPort:
		MapEthernetSwitchPortToDto(entity.(domain.EthernetSwitchPort), dto.(*dtos.EthernetSwitchPortDto))
	//HTTPLog
	case domain.HTTPLog:
		MapHTTPLogEntityToDto(entity.(domain.HTTPLog), dto.(*dtos.HTTPLogDto))
	//AppLog
	case domain.AppLog:
		MapAppLogEntityToDto(entity.(domain.AppLog), dto.(*dtos.AppLogDto))
	//DeviceTemplate
	case domain.DeviceTemplate:
		MapDeviceTemplateToDto(entity.(domain.DeviceTemplate), dto.(*dtos.DeviceTemplateDto))
	//HostNetworkVlan
	case domain.HostNetworkVlan:
		MapHostNetworkVlanToDto(entity.(domain.HostNetworkVlan), dto.(*dtos.HostNetworkVlanDto))
	//TFTPConfig
	case domain.TFTPConfig:
		MapTFTPConfigToDto(entity.(domain.TFTPConfig), dto.(*dtos.TFTPConfigDto))
	//TFTPPathRatio
	case domain.TFTPPathRatio:
		MapTFTPPathToDto(entity.(domain.TFTPPathRatio), dto.(*dtos.TFTPPathDto))
	default:
		return errors.Internal.Newf("can't find route for map entity %+v to dto %+v", dto, entity)
	}
	return nil
}
