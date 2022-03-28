package domain

type EthernetSwitchModel int

const (
	UBIQUITY_US24_250W = iota
)

type EthernetSwitch struct {
	Entity
	Name        string
	Serial      string
	SwitchModel EthernetSwitchModel
	Address     string
	Username    string
	Password    string
	Ports       *[]*EthernetSwitchPort
}
