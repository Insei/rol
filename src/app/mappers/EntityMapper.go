package mappers

import (
	"rol/domain"
)

func GetEmptyEntityFromArrayType(entityArr interface{}) interface{} {
	switch entityArr.(type) {
	case *[]*domain.EthernetSwitch:
		return &domain.EthernetSwitch{}
	case *[]*domain.EthernetSwitchPort:
		return &domain.EthernetSwitchPort{}
	}
	return nil
}
