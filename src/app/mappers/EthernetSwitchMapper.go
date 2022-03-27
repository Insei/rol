package mappers

import (
	"rol/domain/entities"
	"rol/domain/enums"
	"rol/dtos"
)

func MapEthernetSwitchUpdateDto(dto *dtos.EthernetSwitchUpdateDto, entity *entities.EthernetSwitch) {
	entity.SwitchModel = (enums.EthernetSwitchModel)(dto.SwitchModel)
	entity.Name = dto.Name
	entity.Address = dto.Address
	entity.Password = dto.Password
	entity.Username = dto.Username
	entity.Serial = dto.Serial
}

func MapEthernetSwitchCreateDto(dto *dtos.EthernetSwitchCreateDto, entity *entities.EthernetSwitch) {
	entity.SwitchModel = (enums.EthernetSwitchModel)(dto.SwitchModel)
	entity.Name = dto.Name
	entity.Address = dto.Address
	entity.Password = dto.Password
	entity.Username = dto.Username
	entity.Serial = dto.Serial
}

func MapEthernetSwitchDto(entity *entities.EthernetSwitch, dto *dtos.EthernetSwitchDto) {
	dto.Id = entity.ID
	dto.Name = entity.Name
	dto.Address = entity.Address
	dto.Username = entity.Username
	dto.SwitchModel = (int)(entity.SwitchModel)
	dto.Serial = entity.Serial
	dto.CreatedAt = entity.CreatedAt
	dto.UpdatedAt = entity.UpdatedAt
}

func MapEthernetSwitchArrayDto(models *[]*entities.EthernetSwitch, dtoses *[]*dtos.EthernetSwitchDto) {
	for i := range *models {
		dto := &dtos.EthernetSwitchDto{}
		MapEthernetSwitchDto((*models)[i], dto)
		*dtoses = append(*dtoses, dto)
	}
}
