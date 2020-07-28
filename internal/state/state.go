package state

import (
	"sync"
)

// TODO: Automatically create UI and GI?

// State ...
type State struct {
	sync.RWMutex
	config  config
	auths   auths
	current currentState
	desired desiredState
}

// NewState ...
func NewState() State {
	return State{
		auths:   newAuths(),
		current: newCurrentState(),
		desired: newDesiredState(),
	}
}
