package state

import "fmt"

// State ...
type State struct {
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

func (s *State) validateStateUpdateFromGateway(u *StateUpdateFromGateway) error {
	// TODO: break this up like updating

	// TODO: validate auth and that gateway is only updating itself (and inherently
	// its own devices) here or in API

	// validate gateway must already exist
	if _, ok := s.Gateways[GatewayID{ExternalGatewayID: u.Gateway.ExternalID}]; !ok {
		return fmt.Errorf("gateway does not exist: %s", u.Gateway.ExternalID)
	}

	// validate gateway features
	for _, fu := range u.Gateway.Features {
		if err := s.validateFeatureUpdateFromGateway(&fu); err != nil {
			return err
		}
	}

	// validate devices
	for _, du := range u.Gateway.Devices {
		// validate ExternalID must be non-empty
		if du.ExternalID == "" {
			return fmt.Errorf("ExternalID must not be empty on device")
		}

		// validate device features
		for _, fu := range du.Features {
			if err := s.validateFeatureUpdateFromGateway(&fu); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *State) validateFeatureUpdateFromGateway(fu *FeatureUpdateFromGateway) error {
	// valdidate IDReferredToByGateway must be non-empty
	if fu.ExternalID == "" {
		return fmt.Errorf("ExternalID must not be empty on feature")
	}
	return nil
}

func (s *State) updateWithStateUpdateFromGateway(u *StateUpdateFromGateway) error {
	// update gateway and its features
	s.updateWithGatewayUpdateFromGateway(&u.Gateway)

	// update devices and their features
	for _, du := range u.Gateway.Devices {
		s.updateWithDeviceUpdateFromGateway(u.Gateway.ExternalID, &du)
	}

	return nil
}

func (s *State) updateWithGatewayUpdateFromGateway(gu *GatewayUpdateFromGateway) {
	gid := GatewayID{
		ExternalGatewayID: gu.ExternalID,
	}

	// get a copy of the gateway
	g, ok := s.Gateways[gid]
	if !ok {
		// this should never happen becaue we validated!
		// but we could still handle it better than this.
		panic("gateway not found during update")
	}

	// update fields on our gateway
	g.ID = gid
	if gu.Manufacturer != nil {
		g.Manufacturer = *gu.Manufacturer
	}
	if gu.Model != nil {
		g.Model = *gu.Model
	}
	if gu.SerialNumber != nil {
		g.SerialNumber = *gu.SerialNumber
	}

	// replace gateway on state with updated copy
	s.Gateways[gid] = g

	// update the gateway's features
	for _, fu := range gu.Features {
		s.updateWithFeatureUpdateFromGateway(gid.ExternalGatewayID, "", &fu)
	}
}

func (s *State) updateWithDeviceUpdateFromGateway(externalGatewayID string, du *DeviceUpdateFromGateway) {
	did := DeviceID{
		ExternalGatewayID: externalGatewayID,
		ExternalDeviceID:  du.ExternalID,
	}

	// get a copy of the device if we already have it (or create a new one)
	d, ok := s.Devices[did]
	if !ok {
		d = Device{}
	}

	// update fields on our copy
	d.ID = did
	if du.Manufacturer != nil {
		d.Manufacturer = *du.Manufacturer
	}
	if du.Model != nil {
		d.Model = *du.Model
	}
	if du.SerialNumber != nil {
		d.SerialNumber = *du.SerialNumber
	}

	// replace device on state with updated copy
	s.Devices[did] = d

	// update the device's features
	for _, fu := range du.Features {
		s.updateWithFeatureUpdateFromGateway(did.ExternalGatewayID, did.ExternalDeviceID, &fu)
	}
}

func (s *State) updateWithFeatureUpdateFromGateway(externalGatewayID string, externalDeviceID string, fu *FeatureUpdateFromGateway) {
	fid := FeatureID{
		ExternalGatewayID: externalGatewayID,
		ExternalDeviceID:  externalDeviceID,
		ExternalFeatureID: fu.ExternalID,
	}

	// get a copy of the feature if we already have it (or create a new one)
	f, ok := s.Features[fid]
	if !ok {
		f = Feature{}
	}

	// update fields on our copy
	f.ID = fid

	// replace feature on state with updated copy
	s.Features[fid] = f
}
