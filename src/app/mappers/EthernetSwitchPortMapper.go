package mappers

import (
	"rol/domain"
	"rol/dtos"
)

func MapEthernetSwitchPortCreateDto(dto *dtos.EthernetSwitchPortCreateDto, entity *domain.EthernetSwitchPort) {
	entity.Name = dto.Name
	entity.EthernetSwitchId = dto.EthernetSwitchId
	entity.PoeType = (domain.EthernetSwitchPoePortType)(dto.PoeType)
}

func MapEthernetSwitchPortUpdateDto(dto *dtos.EthernetSwitchPortUpdateDto, entity *domain.EthernetSwitchPort) {
	entity.Name = dto.Name
	entity.PoeType = (domain.EthernetSwitchPoePortType)(dto.PoeType)
}

func MapEthernetSwitchPortDto(entity *domain.EthernetSwitchPort, dto *dtos.EthernetSwitchPortDto) {
	dto.Id = entity.ID
	dto.Name = entity.Name
	dto.PoeType = (int)(entity.PoeType)
	dto.EthernetSwitchId = entity.EthernetSwitchId
	dto.CreatedAt = entity.CreatedAt
	dto.UpdatedAt = entity.UpdatedAt
}

func MapEthernetSwitchPortArrayDto(entities *[]*domain.EthernetSwitchPort, dtoses *[]*dtos.EthernetSwitchPortDto) {
	for i := range *entities {
		dto := &dtos.EthernetSwitchPortDto{}
		MapEthernetSwitchPortDto((*entities)[i], dto)
		*dtoses = append(*dtoses, dto)
	}
}
