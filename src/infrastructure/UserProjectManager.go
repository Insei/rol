package infrastructure

import (
	"github.com/sirupsen/logrus"
	"net"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/domain"
)

//UserProjectManager is an object for create all project dependencies such as traffic rule and bridge
type UserProjectManager struct {
	repository     interfaces.IGenericRepository[domain.Project]
	networkManager interfaces.IHostNetworkManager
	logger         *logrus.Logger
}

//NewUserProjectManager UserProjectManager constructor
func NewUserProjectManager(repo interfaces.IGenericRepository[domain.Project], networkManager interfaces.IHostNetworkManager, log *logrus.Logger) interfaces.IProjectManager {
	return UserProjectManager{
		repository:     repo,
		networkManager: networkManager,
		logger:         log,
	}
}

//CreateProject creates project and bridge with postrouting traffic rule for it
//
//Params:
//	project - project entity
//	projectSubnet - subnet for project
//Return:
//	domain.Project - created project
//	error - if an error occurs, otherwise nil
func (p UserProjectManager) CreateProject(project domain.Project, projectSubnet string) (domain.Project, error) {
	bridgeName, err := p.networkManager.CreateBridge(project.Name)
	if err != nil {
		return domain.Project{}, errors.Internal.Wrap(err, "network manager failed to create bridge")
	}
	project.BridgeName = bridgeName

	bridgeIp := projectSubnet[:len(projectSubnet)-1] + "1"
	ip, ipNet, err := net.ParseCIDR(bridgeIp + "/24")
	if err != nil {
		return domain.Project{}, errors.Internal.Wrap(err, "failed to parse cidr")
	}
	ipNet.IP = ip

	trafficRule := domain.HostNetworkTrafficRule{
		Chain:       "POSTROUTING",
		Target:      "MASQUERADE",
		Source:      projectSubnet + "/24",
		Destination: "",
	}

	err = p.networkManager.AddrAdd(bridgeName, *ipNet)
	if err != nil {
		return domain.Project{}, errors.Internal.Wrap(err, "network manager failed to add address to bridge")
	}
	_, err = p.networkManager.CreateTrafficRule("nat", trafficRule)
	if err != nil {
		return domain.Project{}, errors.Internal.Wrap(err, "network manager failed to create traffic rule")
	}
	err = p.networkManager.SetLinkMaster("enp0s8", bridgeName)
	if err != nil {
		return domain.Project{}, errors.Internal.Wrap(err, "network manager failed to set link master")
	}
	err = p.networkManager.SetLinkUp(bridgeName)
	if err != nil {
		return domain.Project{}, errors.Internal.Wrap(err, "network manager failed to set link up")
	}
	err = p.networkManager.SaveConfiguration()
	if err != nil {
		return domain.Project{}, errors.Internal.Wrap(err, "network manager failed to save configuration")
	}

	return project, nil
}

//DeleteProject delete project and its bridge with postrouting traffic rule
//
//Params:
//	project - project entity
//Return:
//	error - if an error occurs, otherwise nil
func (p UserProjectManager) DeleteProject(project domain.Project) error {
	err := p.networkManager.DeleteLinkByName(project.BridgeName)
	if err != nil {
		return errors.Internal.Wrap(err, "network manager failed to delete link by name")
	}

	trafficRule := domain.HostNetworkTrafficRule{
		Chain:       "POSTROUTING",
		Target:      "MASQUERADE",
		Source:      project.Subnet + "/24",
		Destination: "",
	}

	err = p.networkManager.DeleteTrafficRule("nat", trafficRule)
	if err != nil {
		return errors.Internal.Wrap(err, "network manager failed to delete traffic rule")
	}
	err = p.networkManager.SaveConfiguration()
	if err != nil {
		return errors.Internal.Wrap(err, "network manager failed to save configuration")
	}
	return nil
}
