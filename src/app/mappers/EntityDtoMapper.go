package mappers

import (
	"rol/app/interfaces"
	"rol/domain"
	"rol/dtos"
)

func GetEmptyEntity(dto interfaces.IEntityDtoModel) interfaces.IEntityModel {
	switch dto.(type) {
	case *dtos.EthernetSwitchCreateDto, *dtos.EthernetSwitchUpdateDto, *dtos.EthernetSwitchDto:
		return &domain.EthernetSwitch{}
	case *dtos.EthernetSwitchPortCreateDto, *dtos.EthernetSwitchPortUpdateDto, *dtos.EthernetSwitchPortDto:
		return &domain.EthernetSwitchPort{}
	default:
		return nil
	}
}

func GetEntityEmptyArray(dto interface{}) interface{} {
	switch dto.(type) {
	case *[]*dtos.EthernetSwitchDto:
		return &[]*domain.EthernetSwitch{}
	case *[]*dtos.EthernetSwitchPortDto:
		return &[]*domain.EthernetSwitchPort{}
	}
	return nil
}

func Map(source interface{}, dest interface{}) {
	switch source.(type) {
	// EthernetSwitch
	case *dtos.EthernetSwitchCreateDto:
		MapEthernetSwitchCreateDto(source.(*dtos.EthernetSwitchCreateDto), dest.(*domain.EthernetSwitch))
	case *dtos.EthernetSwitchUpdateDto:
		MapEthernetSwitchUpdateDto(source.(*dtos.EthernetSwitchUpdateDto), dest.(*domain.EthernetSwitch))
	case *domain.EthernetSwitch:
		MapEthernetSwitchDto(source.(*domain.EthernetSwitch), dest.(*dtos.EthernetSwitchDto))
	case *[]*domain.EthernetSwitch:
		MapEthernetSwitchArrayDto(source.(*[]*domain.EthernetSwitch), dest.(*[]*dtos.EthernetSwitchDto))
	//EthernetSwitchPort
	case *dtos.EthernetSwitchPortCreateDto:
		MapEthernetSwitchPortCreateDto(source.(*dtos.EthernetSwitchPortCreateDto), dest.(*domain.EthernetSwitchPort))
	case *dtos.EthernetSwitchPortUpdateDto:
		MapEthernetSwitchPortUpdateDto(source.(*dtos.EthernetSwitchPortUpdateDto), dest.(*domain.EthernetSwitchPort))
	case *domain.EthernetSwitchPort:
		MapEthernetSwitchPortDto(source.(*domain.EthernetSwitchPort), dest.(*dtos.EthernetSwitchPortDto))
	case *[]*domain.EthernetSwitchPort:
		MapEthernetSwitchPortArrayDto(source.(*[]*domain.EthernetSwitchPort), dest.(*[]*dtos.EthernetSwitchPortDto))
	default:
		return
	}
}
