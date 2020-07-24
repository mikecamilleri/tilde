package state

import "encoding/json"

// UpdateFromGateway ...
type UpdateFromGateway struct {
	Gateway GatewayUpdateFromGateway
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
