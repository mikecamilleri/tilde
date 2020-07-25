package state

import (
	"encoding/json"
	"fmt"
)

// UpdateFromGateway ...
type UpdateFromGateway struct {
	validated bool
	auth      GatewayAuth
	data      struct {
		Gateway gatewayUpdateFromGateway
	}
}

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

func (u *UpdateFromGateway) validate(c *Current) error {
	// validate authorization (not authentication!) and that gateway is only
	// updating itself (and inherently its own devices due to id construction)
	if u.auth.ID.ExternalGatewayID != u.data.Gateway.ExternalID {
		return fmt.Errorf("gateway attempting to update other than iteself")
	}

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

func (u *UpdateFromGateway) apply(c *Current) error {
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
	Devices      []deviceUpdateFromGateway
	Features     []featureUpdateFromGateway
}

func (gu *gatewayUpdateFromGateway) validate(c *Current) error {
	// validate gateway must already exist
	if _, ok := c.Gateways[GatewayID{ExternalGatewayID: gu.ExternalID}]; !ok {
		return fmt.Errorf("gateway does not exist: %s", gu.ExternalID)
	}

	// validate gateway features
	for _, fu := range gu.Features {
		if err := fu.validate(c); err != nil {
			return err
		}
	}

	return nil
}

func (gu *gatewayUpdateFromGateway) apply(c *Current) error {
	gid := GatewayID{
		ExternalGatewayID: gu.ExternalID,
	}

	// get a copy of the gateway
	g, ok := c.Gateways[gid]
	if !ok {
		return fmt.Errorf("gateway not found during update (this shouldn't happen becaue we validated!)")
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
	c.Gateways[gid] = g

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
	Features     []featureUpdateFromGateway
}

func (du *deviceUpdateFromGateway) validate(c *Current, externalGatewayID string) error {
	// validate ExternalID must be non-empty
	if du.ExternalID == "" {
		return fmt.Errorf("ExternalID must not be empty on device")
	}

	// validate device features
	for _, fu := range du.Features {
		if err := fu.validate(c); err != nil {
			return err
		}
	}

	return nil
}

func (du *deviceUpdateFromGateway) apply(c *Current, externalGatewayID string) {
	did := DeviceID{
		ExternalGatewayID: externalGatewayID,
		ExternalDeviceID:  du.ExternalID,
	}

	// get a copy of the device if we already have it (or create a new one)
	d, ok := c.Devices[did]
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
	c.Devices[did] = d

	// update the device's features
	for _, fu := range du.Features {
		fu.apply(c, did.ExternalGatewayID, did.ExternalDeviceID)
	}
}

type featureUpdateFromGateway struct {
	ExternalID string `json:"id"`
	// Standard    *string
	// user settable?
	// ...
}

func (fu *featureUpdateFromGateway) validate(c *Current) error {
	// valdidate ExternalID must be non-empty
	if fu.ExternalID == "" {
		return fmt.Errorf("ExternalID must not be empty on feature")
	}

	return nil
}

func (fu *featureUpdateFromGateway) apply(c *Current, externalGatewayID string, externalDeviceID string) {
	fid := FeatureID{
		ExternalGatewayID: externalGatewayID,
		ExternalDeviceID:  externalDeviceID,
		ExternalFeatureID: fu.ExternalID,
	}

	// get a copy of the feature if we already have it (or create a new one)
	f, ok := c.Features[fid]
	if !ok {
		f = Feature{}
	}

	// update fields on our copy
	f.ID = fid

	// replace feature on state with updated copy
	c.Features[fid] = f
}
