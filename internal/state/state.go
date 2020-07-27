package state

import (
	"sync"
)

// State ...
// TODO: Rename-- tilde? db?
// TODO: Locks!
type State struct {
	sync.RWMutex
	config  config
	auths   auths
	current currentState
	desired desiredState
}

// NewState ...
func NewState() State {
	return State{
		auths:   newAuths(),
		current: newCurrentState(),
		desired: newDesiredState(),
	}
}

// GatewayAuths ...
func (s *State) GatewayAuths() []GatewayAuth {
	gas := []GatewayAuth{}
	for _, gap := range s.auths.Gateways.all() {
		gas = append(gas, *gap)
	}
	return gas
}

// GatewayAuth ...
func (s *State) GatewayAuth(id GatewayID) (GatewayAuth, error) {
	gap, err := s.auths.Gateways.get(id)
	return *gap, err
}

// NewGatewayAuth ...
func (s *State) NewGatewayAuth(id GatewayID, username string, password string) (GatewayAuth, error) {
	gap, err := s.auths.Gateways.create(id, username, password)
	return *gap, err
}

// DeleteGatewayAuth ...
func (s *State) DeleteGatewayAuth(id GatewayID) {
	s.auths.Gateways.delete(id)
}

// GatewayOTPs ...
func (s *State) GatewayOTPs() []GatewayOTP {
	gos := []GatewayOTP{}
	for _, gop := range s.auths.GatewayOTPs.all() {
		gos = append(gos, *gop)
	}
	return gos
}

// GatewayOTP ...
func (s *State) GatewayOTP(id GatewayID) (GatewayOTP, error) {
	gop, err := s.auths.GatewayOTPs.get(id)
	return *gop, err
}

// NewGatewayOTP ...
func (s *State) NewGatewayOTP(id GatewayID, username string, password string) (GatewayOTP, error) {
	gop, err := s.auths.GatewayOTPs.create(id, username, password, s.config.GatewayOTPExpirationSeconds)
	return *gop, err
}

// DeleteGatewayOTP ...
func (s *State) DeleteGatewayOTP(id GatewayID) {
	s.auths.GatewayOTPs.delete(id)
}

// UserAuths ...
func (s *State) UserAuths() []UserAuth {
	uas := []UserAuth{}
	for _, uap := range s.auths.Users.all() {
		uas = append(uas, *uap)
	}
	return uas
}

// UserAuth ...
func (s *State) UserAuth(id UserID) (UserAuth, error) {
	uap, err := s.auths.Users.get(id)
	return *uap, err
}

// NewUserAuth ...
func (s *State) NewUserAuth(id UserID, username string, password string) (UserAuth, error) {
	uap, err := s.auths.Users.create(id, username, password)
	return *uap, err
}

// UpdateUserAuthPassword ...
func (s *State) UpdateUserAuthPassword(id UserID, password string) error {
	uap, err := s.auths.Users.get(id)
	if err != nil {
		return err
	}
	uap.setPassword(password)
	return nil
}

// DeleteUserAuth ...
func (s *State) DeleteUserAuth(id UserID) {
	s.auths.Users.delete(id)
}

type config struct {
	GatewayOTPExpirationSeconds int
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

// desiredState ...
type desiredState struct {
	Features map[FeatureID]*Feature
}

// newCurrentState is used to avoid nil pointer problems
func newDesiredState() desiredState {
	return desiredState{
		Features: make(map[FeatureID]*Feature),
	}
}

// TODO:
// - add non-exported functions to create new Gateway/Device/Feature
// - add exported methods to State that allow for addiotn of new
//   Gateway/Device/Feature with validation of relationships. Don't forget
//   locks!

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

// UserID ...
type UserID struct {
	UserID string
}

func (g UserID) String() string {
	return g.UserID
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
