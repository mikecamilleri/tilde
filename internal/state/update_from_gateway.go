package state

import (
	"encoding/json"
	"fmt"
)

// UpdateFromGateway ...
type UpdateFromGateway struct {
	validated bool
	Gateway   GatewayUpdateFromGateway
}

// GatewayUpdateFromGateway ...
type GatewayUpdateFromGateway struct {
	ExternalID   string
	Manufacturer *string
	Model        *string
	SerialNumber *string
	Devices      []DeviceUpdateFromGateway
	Features     []FeatureUpdateFromGateway
}

// DeviceUpdateFromGateway ...
type DeviceUpdateFromGateway struct {
	ExternalID   string
	Manufacturer *string
	Model        *string
	SerialNumber *string
	Features     []FeatureUpdateFromGateway
}

// FeatureUpdateFromGateway ...
type FeatureUpdateFromGateway struct {
	ExternalID string
	// Standard    *string
	// user settable?
	// ...
}

// NewUpdateFromGateway unmarshals the JSON update from the gateway into a
// StateUpdateFromGateway struct. Extra fields are ignored.
func NewUpdateFromGateway(updateJSONBytes []byte) (UpdateFromGateway, error) {
	u := UpdateFromGateway{}
	if err := json.Unmarshal(updateJSONBytes, &u); err != nil {
		return u, err
	}
	return u, nil
}

func (u *UpdateFromGateway) validate(s *State) error {
	// TODO: validate auth and that gateway is only updating itself (and inherently
	// its own devices due to id constrcution) here or in API

	// validate gateway and its features
	if err := u.Gateway.validate(s); err != nil {
		return err
	}

	// validate devices and their features
	for _, du := range u.Gateway.Devices {
		// validate ExternalID must be non-empty
		if err := du.validate(s, u.Gateway.ExternalID); err != nil {
			return err
		}
	}

	// mark as validated
	u.validated = true

	return nil
}

func (gu *GatewayUpdateFromGateway) validate(s *State) error {
	// validate gateway must already exist
	if _, ok := s.Gateways[GatewayID{ExternalGatewayID: gu.ExternalID}]; !ok {
		return fmt.Errorf("gateway does not exist: %s", gu.ExternalID)
	}

	// validate gateway features
	for _, fu := range gu.Features {
		if err := fu.validate(s); err != nil {
			return err
		}
	}

	return nil
}

func (du *DeviceUpdateFromGateway) validate(s *State, externalGatewayID string) error {
	// validate ExternalID must be non-empty
	if du.ExternalID == "" {
		return fmt.Errorf("ExternalID must not be empty on device")
	}

	// validate device features
	for _, fu := range du.Features {
		if err := fu.validate(s); err != nil {
			return err
		}
	}

	return nil
}

func (fu *FeatureUpdateFromGateway) validate(s *State) error {
	// valdidate IDReferredToByGateway must be non-empty
	if fu.ExternalID == "" {
		return fmt.Errorf("ExternalID must not be empty on feature")
	}

	return nil
}

func (u *UpdateFromGateway) apply(s *State) error {
	if !u.validated {
		return fmt.Errorf("update not validated before applying")
	}

	// update gateway and its features
	if err := u.Gateway.apply(s); err != nil {
		return err
	}

	// update devices and their features
	for _, du := range u.Gateway.Devices {
		if err := du.apply(s, u.Gateway.ExternalID); err != nil {
			return err
		}
	}

	return nil
}

func (gu *GatewayUpdateFromGateway) apply(s *State) error {
	gid := GatewayID{
		ExternalGatewayID: gu.ExternalID,
	}

	// get a copy of the gateway
	g, ok := s.Gateways[gid]
	if !ok {
		// this should never happen becaue we validated!
		// but we could still handle it better than this.
		return fmt.Errorf("gateway not found during update (this shouldn't happen!)")
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
		if err := fu.apply(s, gid.ExternalGatewayID, ""); err != nil {
			return err
		}
	}

	return nil
}

func (du *DeviceUpdateFromGateway) apply(s *State, externalGatewayID string) error {
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
		if err := fu.apply(s, did.ExternalGatewayID, did.ExternalDeviceID); err != nil {
			return err
		}
	}

	return nil
}

func (fu *FeatureUpdateFromGateway) apply(s *State, externalGatewayID string, externalDeviceID string) error {
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

	return nil
}
