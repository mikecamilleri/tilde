package state

import "sync"

// State ...
type State struct {
	sync.RWMutex
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
	// Features              []FeatureID
}

// Feature ...
type Feature struct {
	ID          FeatureID
	Name        string
	Description string
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

	if err := u.validate(s); err != nil {
		return err
	}

	if err := u.apply(s); err != nil {
		return err
	}

	return nil
}
