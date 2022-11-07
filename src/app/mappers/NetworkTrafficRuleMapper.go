package mappers

import (
	"github.com/coreos/go-iptables/iptables"
	"rol/domain"
	"rol/dtos"
)

func MapStatToTrafficRule(stat iptables.Stat, rule *domain.HostNetworkTrafficRule) {
	rule.Source = stat.Source.String()
	rule.Destination = stat.Destination.String()
	rule.Target = stat.Target
}

func MapHostNetworkTrafficRuleEntityToDto(entity domain.HostNetworkTrafficRule, dto *dtos.HostNetworkTrafficRuleDto) {
	dto.Chain = entity.Chain
	dto.Target = entity.Target
	dto.Source = entity.Source
	dto.Destination = entity.Destination
}

func MapHostNetworkTrafficRuleCreateDtoToEntity(dto dtos.HostNetworkTrafficRuleCreateDto, entity *domain.HostNetworkTrafficRule) {
	entity.Chain = dto.Chain
	entity.Target = dto.Target
	entity.Source = dto.Source
	entity.Destination = dto.Destination
}

func MapHostNetworkTrafficRuleDeleteDtoToEntity(dto dtos.HostNetworkTrafficRuleDeleteDto, entity *domain.HostNetworkTrafficRule) {
	entity.Chain = dto.Chain
	entity.Target = dto.Target
	entity.Source = dto.Source
	entity.Destination = dto.Destination
}
