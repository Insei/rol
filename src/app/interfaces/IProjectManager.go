package interfaces

import "rol/domain"

//IProjectManager is an interface for project manager
type IProjectManager interface {
	//CreateProject creates project and bridge with postrouting traffic rule for it
	//
	//Params:
	//	project - project entity
	//	projectSubnet - subnet for project
	//Return:
	//	domain.Project - created project
	//	error - if an error occurs, otherwise nil
	CreateProject(project domain.Project, projectSubnet string) (domain.Project, error)
	//DeleteProject delete project and its bridge with postrouting traffic rule
	//
	//Params:
	//	project - project entity
	//Return:
	//	error - if an error occurs, otherwise nil
	DeleteProject(project domain.Project) error
}
