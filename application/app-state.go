package application

import (
	"net"
	"sicobo/tools"
)

/*
 * Global states / configuration access for the current lib app
 */

// State : Global application environment
var State = initState()

// ApplicationState : global state model for application
type ApplicationState struct {
	Config     Configuration
	OutBountIP net.IP
}

// Get preferred outbound ip of this machine
func initState() ApplicationState {
	return ApplicationState{loadConfiguration(), tools.GetOutboundIP()}
}
