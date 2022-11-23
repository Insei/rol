package interfaces

import (
	"net"
	"rol/domain"
)

//IHostNetworkManager is an interface for network manager
type IHostNetworkManager interface {
	//GetList gets list of host network interfaces
	//
	//Return:
	//	[]interfaces.IHostNetworkLink - list of interfaces
	//	error - if an error occurs, otherwise nil
	GetList() ([]IHostNetworkLink, error)
	//GetByName gets host network interface by its name
	//
	//Params:
	//	name - interface name
	//Return:
	//	interfaces.IHostNetworkLink - host network interface
	//	error - if an error occurs, otherwise nil
	GetByName(name string) (IHostNetworkLink, error)
	//CreateVlan creates vlan on host
	//
	//Params:
	//	master - name of the master network interface
	//	vlanID - ID to be set for vlan
	//Return:
	//	string - new vlan name that will be {master}.{vlanID}
	//	error - if an error occurs, otherwise nil
	CreateVlan(master string, vlanID int) (string, error)
	//CreateBridge creates bridge on host
	//
	//Params:
	//	name - new bridge name
	//Return:
	//	string - new bridge name that will be rol.br.{name}
	//	error - if an error occurs, otherwise nil
	CreateBridge(name string) (string, error)
	//SetLinkUp enables the link
	//
	//Params:
	//	linkName - name of the link
	//Return:
	//	error - if an error occurs, otherwise nil
	SetLinkUp(linkName string) error
	//SetLinkMaster set master for link
	//
	//Params:
	//	slaveName - name of link that will be slave
	//	masterName - name of link that will be slave
	//Return:
	//	error - if an error occurs, otherwise nil
	SetLinkMaster(slaveName, masterName string) error
	//UnsetLinkMaster removes the master of the link
	//
	//Params:
	//	linkName - name of the link
	//Return:
	//	error - if an error occurs, otherwise nil
	UnsetLinkMaster(linkName string) error
	//DeleteLinkByName deletes interface on host by its name
	//
	//Params:
	//	name - interface name
	//Return
	//	error - if an error occurs, otherwise nil
	DeleteLinkByName(name string) error
	//AddrAdd Add new ip address for network interface
	//
	//Params:
	//	linkName - name of the interface
	//	addr - ip address with mask net.IPNet
	//Return:
	//	error - if an error occurs, otherwise nil
	AddrAdd(linkName string, addr net.IPNet) error
	//AddrDelete Delete ip address from network interface
	//
	//Params:
	//	linkName - name of the interface
	//	addr - ip address with mask net.IPNet
	//Return:
	//	error - if an error occurs, otherwise nil
	AddrDelete(linkName string, addr net.IPNet) error
	//CreateTrafficRule Create netfilter traffic rule for specified table
	//
	//Params:
	//	table - table to create a rule
	//	rule - rule entity
	//Return:
	//	domain.HostNetworkTrafficRule - new traffic rule
	//	error - if an error occurs, otherwise nil
	CreateTrafficRule(table string, rule domain.HostNetworkTrafficRule) (domain.HostNetworkTrafficRule, error)
	//DeleteTrafficRule Delete netfilter traffic rule in specified table
	//
	//Params:
	//	table - table to delete a rule
	//	rule - rule entity
	//Return:
	//	error - if an error occurs, otherwise nil
	DeleteTrafficRule(table string, rule domain.HostNetworkTrafficRule) error
	//GetChainRules Get selected netfilter chain rules at specified table
	//
	//Params:
	//	table - table to get a rules
	//	chain - chain where we get the rules
	//Return:
	//	[]domain.HostNetworkTrafficRule - slice of rules
	//	error - if an error occurs, otherwise nil
	GetChainRules(table string, chain string) ([]domain.HostNetworkTrafficRule, error)
	//GetTableRules Get specified netfilter table rules
	//
	//Params:
	//	table - table to get a rules
	//Return:
	//	[]domain.HostNetworkTrafficRule - slice of rules
	//	error - if an error occurs, otherwise nil
	GetTableRules(table string) ([]domain.HostNetworkTrafficRule, error)
	//SaveConfiguration save current host network configuration to the configuration storage
	//
	//Return:
	//	error - if an error occurs, otherwise nil
	SaveConfiguration() error
	//RestoreFromBackup restore and apply host network configuration from backup configuration
	//
	//Return:
	//	error - if an error occurs, otherwise nil
	RestoreFromBackup() error
	//ResetChanges Load and apply host network configuration from saved configuration
	//
	//Return:
	//	error - if an error occurs, otherwise nil
	ResetChanges() error
	//HasUnsavedChanges Gets a flag about unsaved changes
	//
	//Return:
	//	bool - if unsaved changes exist - true, otherwise false
	HasUnsavedChanges() bool
}
