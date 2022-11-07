package domain

//HostNetworkTrafficRule netfilter traffic rule
type HostNetworkTrafficRule struct {
	//Chain rule chain
	Chain string
	//Target rule action like ACCEPT, MASQUERADE, DROP, etc.
	Target string
	//Source packets source
	Source string
	//Destination packets destination
	Destination string
}
