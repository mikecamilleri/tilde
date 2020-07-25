package state

import (
	"sync"
	"time"
)

// State ...
type State struct {
	sync.RWMutex
	auths   Auths
	current Current
}

// Auths ...
type Auths struct {
	Gateways    map[GatewayID]GatewayAuth
	GatewayOTPs map[GatewayID]GatewayOTP
	// Users map[UserID]User
}

// GatewayAuth ...
type GatewayAuth struct {
	ID           GatewayID
	Username     string
	passwordHash string
}

// GatewayOTP ...
type GatewayOTP struct {
	ID       GatewayID
	Username string
	Expires  time.Time
	otpHash  string
}

// Current ...
type Current struct {
	Gateways map[GatewayID]Gateway
	Devices  map[DeviceID]Device
	Features map[FeatureID]Feature
}

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
