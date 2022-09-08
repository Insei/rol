package services

import (
	"net"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/mappers"
	"rol/app/utils"
	"rol/app/validators"
	"rol/dtos"
	"strings"
)

const vlanNotFound = "vlan is not exist on the host"
const bridgeNotFound = "bridge is not exist on the host"
const parentNotFound = "parent interface is not exist on the host"
const slaveNotFound = "slave interface is not exist on the host"
const setAddressesFailed = "set addresses to the vlan fail"

//HostNetworkService is a struct for host network vlan service
type HostNetworkService struct {
	manager interfaces.IHostNetworkManager
}

//NewHostNetworkService is a constructor for HostNetworkService
//
//Params:
//	manager - host network manager
//Return:
//	HostNetworkService - instance of network vlan service
func NewHostNetworkService(manager interfaces.IHostNetworkManager) *HostNetworkService {
	return &HostNetworkService{
		manager: manager,
	}
}

//GetVlanList gets list of host vlans
//
//Return:
//	[]dtos.HostNetworkVlanDto - slice of vlan dtos
//	error - if an error occurs, otherwise nil
func (h *HostNetworkService) GetVlanList() ([]dtos.HostNetworkVlanDto, error) {
	out := []dtos.HostNetworkVlanDto{}
	links, err := h.manager.GetList()
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error getting link list")
	}
	for _, link := range links {
		if link.GetType() == "vlan" && strings.Contains(link.GetName(), "rol.") {
			var dto dtos.HostNetworkVlanDto
			err = mappers.MapEntityToDto(link, &dto)
			if err != nil {
				return nil, errors.Internal.Wrap(err, "error mapping vlan")
			}
			out = append(out, dto)
		}
	}
	return out, nil
}

//GetVlanByName gets vlan by name
//
//Params:
//	vlanName - name of the vlan
//Return:
//	dtos.HostNetworkVlanDto - vlan dto
//	error - if an error occurs, otherwise nil
func (h *HostNetworkService) GetVlanByName(vlanName string) (dtos.HostNetworkVlanDto, error) {
	link, err := h.manager.GetByName(vlanName)
	if err != nil {
		return dtos.HostNetworkVlanDto{}, errors.Internal.Wrap(err, "error getting vlan by name")
	}
	if link == nil || link.GetType() != "vlan" || !strings.Contains(link.GetName(), "rol.") {
		return dtos.HostNetworkVlanDto{}, errors.NotFound.New(vlanNotFound)
	}
	var dto dtos.HostNetworkVlanDto
	err = mappers.MapEntityToDto(link, &dto)
	if err != nil {
		return dtos.HostNetworkVlanDto{}, errors.Internal.Wrap(err, "error mapping vlan")
	}
	return dto, nil
}

//CreateVlan new vlan on host
//
//Params:
//	vlan - vlan create dto
//Return:
//	dtos.HostNetworkVlanDto - created host network vlan
//	error - if an error occurs, otherwise nil
func (h *HostNetworkService) CreateVlan(createDto dtos.HostNetworkVlanCreateDto) (dtos.HostNetworkVlanDto, error) {
	dto := dtos.HostNetworkVlanDto{}
	err := validators.ValidateHostNetworkVlanCreateDto(createDto)
	if err != nil {
		return dto, err
	}
	_, err = h.manager.GetByName(createDto.Parent)
	if err != nil {
		if !errors.As(err, errors.NotFound) {
			return dto, errors.Internal.Wrap(err, "failed to check existence of master vlan interface")
		}
		err1 := errors.Validation.New(errors.ValidationErrorMessage)
		return dto, errors.AddErrorContext(err1, "Parent", parentNotFound)
	}
	vlanName, err := h.manager.CreateVlan(createDto.Parent, createDto.VlanID)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error creating vlan")
	}
	err = h.manager.SetLinkUp(vlanName)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "set vlan up failed")
	}
	link, err := h.manager.GetByName(vlanName)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error getting vlan by name")
	}
	err = h.syncAddresses(link, createDto.Addresses)
	if err != nil {
		err1 := h.manager.ResetChanges()
		if err1 != nil {
			return dto, errors.Internal.Wrap(err1, "fatal: failed to reset changes after fail with setup address")
		}
		return dto, err
	}
	//Update link from manager
	link, err = h.manager.GetByName(vlanName)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error getting vlan by name")
	}
	err = mappers.MapEntityToDto(link, &dto)
	if err != nil {
		return dtos.HostNetworkVlanDto{}, errors.Internal.Wrap(err, "error mapping vlan")
	}
	return dto, nil
}

func (h *HostNetworkService) syncAddresses(link interfaces.IHostNetworkLink, addresses []string) error {
	currAddresses := link.GetAddresses()
	linkName := link.GetName()
	currAddressStrSlice := []string{}
	for _, address := range currAddresses {
		currAddressStrSlice = append(currAddressStrSlice, address.String())
	}
	deletedCidrSlice, addedCidrSlice := utils.SliceDiffElements(currAddressStrSlice, addresses)
	for _, deletedCidr := range deletedCidrSlice {
		ip, address, err := net.ParseCIDR(deletedCidr)
		if err != nil {
			return errors.Internal.New("failed to parse CIDR")
		}
		address.IP = ip
		err = h.manager.AddrDelete(linkName, *address)
		if err != nil {
			err1 := h.manager.ResetChanges()
			if err1 != nil {
				return errors.Internal.Wrap(err, "fatal: failed to reset changes after fail with setup address")
			}
			return errors.Internal.Wrap(err, setAddressesFailed)
		}
	}
	for _, addedCidr := range addedCidrSlice {
		ip, address, err := net.ParseCIDR(addedCidr)
		if err != nil {
			return errors.Internal.New("failed to parse CIDR")
		}
		address.IP = ip
		err = h.manager.AddrAdd(linkName, *address)
		if err != nil {
			err1 := h.manager.ResetChanges()
			if err1 != nil {
				return errors.Internal.Wrap(err, "fatal: failed to reset changes after fail with setup address")
			}
			return errors.Internal.Wrap(err, setAddressesFailed)
		}
	}
	return nil
}

//UpdateVlan vlan on host
//
//Params:
//	vlanName - vlan name
//	updateDto - vlan update dto
//Return:
//	dtos.HostNetworkVlanDto - updated host network vlan
//	error - if an error occurs, otherwise nil
func (h *HostNetworkService) UpdateVlan(vlanName string, updateDto dtos.HostNetworkVlanUpdateDto) (dtos.HostNetworkVlanDto, error) {
	dto := dtos.HostNetworkVlanDto{}
	err := validators.ValidateHostNetworkVlanUpdateDto(updateDto)
	if err != nil {
		return dto, err
	}
	link, err := h.manager.GetByName(vlanName)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error getting vlan by name")
	}
	if link == nil || (link.GetType() != "vlan" && strings.Contains(link.GetName(), "rol.")) {
		return dto, errors.NotFound.New("vlan not found")
	}
	err = h.syncAddresses(link, updateDto.Addresses)
	if err != nil {
		return dto, err
	}
	//Update link from manager
	link, err = h.manager.GetByName(vlanName)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error getting vlan by name")
	}
	err = mappers.MapEntityToDto(link, &dto)
	if err != nil {
		return dtos.HostNetworkVlanDto{}, errors.Internal.Wrap(err, "error mapping vlan")
	}
	return dto, nil
}

//DeleteVlan deletes vlan on host by its name
//
//Params:
//	vlanName - vlan name
//Return
//	error - if an error occurs, otherwise nil
func (h *HostNetworkService) DeleteVlan(vlanName string) error {
	if !strings.Contains(vlanName, "rol.") {
		return errors.NotFound.New(vlanNotFound)
	}
	link, err := h.manager.GetByName(vlanName)
	if err != nil {
		if !errors.As(err, errors.NotFound) {
			return errors.Internal.Wrap(err, "failed to check existence of vlan interface")
		}
		return err
	}
	if link == nil || link.GetType() != "vlan" {
		return errors.NotFound.New(vlanNotFound)
	}
	err = h.manager.DeleteLinkByName(vlanName)
	if err != nil {
		return errors.Internal.Wrap(err, "delete vlan failed")
	}
	return nil
}
