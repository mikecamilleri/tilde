package state

import (
	"sync"
	"time"
)

// State ...
// TODO: rename? CoreState? SystemState?
type State struct {
	sync.RWMutex
	auths   Auths
	current CurrentState // TODO: rename?
}

// NewState ...
func NewState() State {
	return State{
		auths:   newAuths(),
		current: newCurrentState(),
	}
}

// Auths ...
type Auths struct {
	Gateways    map[GatewayID]*GatewayAuth
	GatewayOTPs map[GatewayID]*GatewayOTP
	// Users map[UserID]User
}

// newAuths ...
func newAuths() Auths {
	return Auths{
		Gateways:    make(map[GatewayID]*GatewayAuth),
		GatewayOTPs: make(map[GatewayID]*GatewayOTP),
	}
}

// TODO:
// - add non-exported functions to create new GatewayAuth/GatewayOTP
// - add exported methods to State that allow for addition of new
//   GatewayAuth/GatewayOTP with validation of relationships. Don't forget
//   locks!

// GatewayAuth ...
type GatewayAuth struct {
	ID           GatewayID
	Username     string
	passwordHash string
}

// GatewayOTP ...
type GatewayOTP struct {
	ID       GatewayID // future gateway id? string instead?
	Username string
	Expires  time.Time
	otpHash  string
}

// CurrentState ...
type CurrentState struct {
	Gateways map[GatewayID]*Gateway
	Devices  map[DeviceID]*Device
	Features map[FeatureID]*Feature
}

// NewCurrentState ...
func newCurrentState() CurrentState {
	return CurrentState{
		Devices:  make(map[DeviceID]*Device),
		Gateways: make(map[GatewayID]*Gateway),
		Features: make(map[FeatureID]*Feature),
	}
}

// TODO:
// - add non-exported functions to create new Gateway/Device/Feature
// - add exported methods to State that allow for addiotn of new
//   Gateway/Device/Feature with validation of relationships. Don't forget
//   locks!

// Gateway ...
type Gateway struct {
	ID           GatewayID
	Name         string
	Description  string
	Manufacturer string
	Model        string
	SerialNumber string
	Active       bool
	Ignore       bool
	// Features     []FeatureID
	// Devices      []DeviceID
}

// Device ...
type Device struct {
	ID           DeviceID
	Name         string
	Description  string
	Manufacturer string
	Model        string
	SerialNumber string
	Active       bool
	Ignore       bool
	// Features              []FeatureID
}

// Feature ...
type Feature struct {
	ID          FeatureID
	Name        string
	Description string
	Active      bool
	Ignore      bool
	// Standard    string
	// user settable?
	// ...
}

// GatewayID ...
type GatewayID struct {
	ExternalGatewayID string
}

func (g GatewayID) String() string {
	return g.ExternalGatewayID
}

// DeviceID ...
type DeviceID struct {
	ExternalGatewayID string
	ExternalDeviceID  string
}

func (d DeviceID) String() string {
	return d.ExternalGatewayID + "-" + d.ExternalDeviceID
}

// FeatureID ...
type FeatureID struct {
	ExternalGatewayID string
	ExternalDeviceID  string // empty for gateway features
	ExternalFeatureID string
}

func (f FeatureID) String() string {
	return f.ExternalGatewayID + "-" + f.ExternalDeviceID + "-" + f.ExternalFeatureID
}

// ApplyUpdateFromGateway ...
func (s *State) ApplyUpdateFromGateway(u UpdateFromGateway) error {
	s.Lock()
	defer s.Unlock()

	if err := u.validate(&s.current); err != nil {
		return err
	}

	if err := u.apply(&s.current); err != nil {
		return err
	}

	return nil
}
