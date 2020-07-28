package state

// UserID ...
type UserID struct {
	UserID string
}

// GatewayID ...
type GatewayID struct {
	ExternalGatewayID string
}

// DeviceID ...
type DeviceID struct {
	ExternalGatewayID string
	ExternalDeviceID  string
}

// FeatureID ...
type FeatureID struct {
	ExternalGatewayID string
	ExternalDeviceID  string // empty for gateway features
	ExternalFeatureID string
}
