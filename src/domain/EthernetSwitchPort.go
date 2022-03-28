package domain

type EthernetSwitchPoePortType int

const (
	POE = iota + 1
	POE_PLUS
	PASSIVE_V24
	NONE = 10
)

type EthernetSwitchPort struct {
	Entity
	Name             string
	EthernetSwitchId uint
	PoeType          EthernetSwitchPoePortType
}
