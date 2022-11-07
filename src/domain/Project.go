package domain

//Project entity of a project
type Project struct {
	Entity
	//Name project name
	Name string
	//BridgeName name of a bridge that related to project
	BridgeName string
	//Subnet project subnet
	Subnet string
}
