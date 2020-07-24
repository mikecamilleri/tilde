package state

import "encoding/json"

// StateUpdateFromGateway ...
type StateUpdateFromGateway struct {
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

// NewStateUpdateFromGateway unmarshals the JSON update from the gateway into a
// StateUpdateFromGateway struct. Extra fields are ignored.
func NewStateUpdateFromGateway(updateJSONBytes []byte) (StateUpdateFromGateway, error) {
	u := StateUpdateFromGateway{}
	if err := json.Unmarshal(updateJSONBytes, &u); err != nil {
		return u, err
	}
	return u, nil
}
