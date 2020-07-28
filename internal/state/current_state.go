package state

// TODO: wrap maps in types as in auths

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
	// Setting(s)  ?
	// Reading(s)  ?
	// user settable?
	// ...
}

type currentState struct {
	Gateways map[GatewayID]*Gateway
	Devices  map[DeviceID]*Device
	Features map[FeatureID]*Feature
}

func newCurrentState() currentState {
	return currentState{
		Devices:  make(map[DeviceID]*Device),
		Gateways: make(map[GatewayID]*Gateway),
		Features: make(map[FeatureID]*Feature),
	}
}
