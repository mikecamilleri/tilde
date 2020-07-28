package state

import (
	"encoding/json"
	"fmt"
)

// GatewayInterface ...
type GatewayInterface struct {
	state *State
}

// NewGatewayInterface ...
func NewGatewayInterface(state *State) GatewayInterface {
	return GatewayInterface{state: state}
}

// ApplyUpdateFromGateway ...
func (gu *GatewayInterface) ApplyUpdateFromGateway(u UpdateFromGateway) error {
	s := gu.state
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

// UpdateFromGateway ...
type UpdateFromGateway struct {
	validated bool
	auth      GatewayAuth
	data      struct {
		Gateway *gatewayUpdateFromGateway
	}
}

// TODO: Locks and consider making methods of UI

// NewUpdateFromGateway unmarshals the JSON update from the gateway into an
// UpdateFromGateway struct. Extra fields are ignored.
func NewUpdateFromGateway(auth GatewayAuth, updateJSONBytes []byte) (UpdateFromGateway, error) {
	u := UpdateFromGateway{
		auth: auth,
	}
	if err := json.Unmarshal(updateJSONBytes, &u.data); err != nil {
		return u, err
	}
	return u, nil
}

func (u *UpdateFromGateway) authenticate(c *currentState) error {
	// validate authentication
	// TODO! ...

	// validate authorization. A gateway may only update itself (and inherently
	// its own devices due to id construction)
	if u.auth.ID.ExternalGatewayID != u.data.Gateway.ExternalID {
		return fmt.Errorf("not authorized: gateway attempting to update other than iteself")
	}

	return nil
}

func (u *UpdateFromGateway) validate(c *currentState) error {
	// validate gateway and its features
	if err := u.data.Gateway.validate(c); err != nil {
		return err
	}

	// validate devices and their features
	for _, du := range u.data.Gateway.Devices {
		// validate ExternalID must be non-empty
		if err := du.validate(c, u.data.Gateway.ExternalID); err != nil {
			return err
		}
	}

	// mark as validated
	u.validated = true

	return nil
}

func (u *UpdateFromGateway) apply(c *currentState) error {
	// is this validated?
	if !u.validated {
		return fmt.Errorf("update not validated before applying")
	}

	// update gateway and its features
	if err := u.data.Gateway.apply(c); err != nil {
		return err
	}

	// update devices and their features
	for _, du := range u.data.Gateway.Devices {
		du.apply(c, u.data.Gateway.ExternalID)
	}

	return nil
}

type gatewayUpdateFromGateway struct {
	ExternalID   string `json:"id"`
	Manufacturer *string
	Model        *string
	SerialNumber *string
	Active       *bool
	Devices      []*deviceUpdateFromGateway
	Features     []*featureUpdateFromGateway
}

func (gu *gatewayUpdateFromGateway) validate(c *currentState) error {
	// we are receiving an update, so mark as active unless spcified otherwise
	// in update
	if gu.Active == nil {
		gu.Active = new(bool)
		*gu.Active = true
	}

	// validate gateway features
	for _, fu := range gu.Features {
		if err := fu.validate(c); err != nil {
			return err
		}
	}

	return nil
}

func (gu *gatewayUpdateFromGateway) apply(c *currentState) error {
	gid := GatewayID{
		ExternalGatewayID: gu.ExternalID,
	}

	// get a refernce to the device if we already have it, or create a new one
	g, ok := c.Gateways[gid]
	if !ok {
		c.Gateways[gid] = new(Gateway)
		c.Gateways[gid].ID = gid
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
	if gu.Active != nil {
		g.Active = *gu.Active
	}

	// update the gateway's features
	for _, fu := range gu.Features {
		fu.apply(c, gid.ExternalGatewayID, "")
	}

	return nil
}

type deviceUpdateFromGateway struct {
	ExternalID   string `json:"id"`
	Manufacturer *string
	Model        *string
	SerialNumber *string
	Active       *bool
	Features     []*featureUpdateFromGateway
}

func (du *deviceUpdateFromGateway) validate(c *currentState, externalGatewayID string) error {
	// validate ExternalID must be non-empty
	if du.ExternalID == "" {
		return fmt.Errorf("ExternalID must not be empty on device")
	}

	// we are receiving an update, so mark as active unless spcified otherwise
	// in update
	if du.Active == nil {
		du.Active = new(bool)
		*du.Active = true
	}

	// validate device features
	for _, fu := range du.Features {
		if err := fu.validate(c); err != nil {
			return err
		}
	}

	return nil
}

func (du *deviceUpdateFromGateway) apply(c *currentState, externalGatewayID string) {
	did := DeviceID{
		ExternalGatewayID: externalGatewayID,
		ExternalDeviceID:  du.ExternalID,
	}

	// get a refernce to the device if we already have it, or create a new one
	d, ok := c.Devices[did]
	if !ok {
		c.Devices[did] = new(Device)
		c.Devices[did].ID = did
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
	if du.Active != nil {
		d.Active = *du.Active
	}

	// update the device's features
	for _, fu := range du.Features {
		fu.apply(c, did.ExternalGatewayID, did.ExternalDeviceID)
	}
}

type featureUpdateFromGateway struct {
	ExternalID string `json:"id"`
	Active     *bool
	// Standard    *string
	// user settable?
	// ...
}

func (fu *featureUpdateFromGateway) validate(c *currentState) error {
	// valdidate ExternalID must be non-empty
	if fu.ExternalID == "" {
		return fmt.Errorf("ExternalID must not be empty on feature")
	}

	// we are receiving an update, so mark as active unless spcified otherwise
	// in update
	if fu.Active == nil {
		fu.Active = new(bool)
		*fu.Active = true
	}

	return nil
}

func (fu *featureUpdateFromGateway) apply(c *currentState, externalGatewayID string, externalDeviceID string) {
	fid := FeatureID{
		ExternalGatewayID: externalGatewayID,
		ExternalDeviceID:  externalDeviceID,
		ExternalFeatureID: fu.ExternalID,
	}

	// get a refernce to the device if we already have it, or create a new one
	f, ok := c.Features[fid]
	if !ok {
		c.Features[fid] = new(Feature)
		c.Features[fid].ID = fid
	}

	// update fields on our copy
	if fu.Active != nil {
		f.Active = *fu.Active
	}

	// update fields on our copy
	f.ID = fid
}
