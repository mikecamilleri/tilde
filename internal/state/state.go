package state

import (
	"fmt"
	"sync"
	"time"
)

// State ...
// TODO: Rename-- tilde? db?
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
	auths := []GatewayAuth{}
	for _, v := range s.auths.Gateways {
		auths = append(auths, *v)
	}
	return auths
}

// GatewayAuth ...
func (s *State) GatewayAuth(id GatewayID) (GatewayAuth, error) {
	auth, ok := s.auths.Gateways[id]
	if !ok {
		return *auth, fmt.Errorf("auth not found")
	}
	return *auth, nil
}

// NewGatewayAuth ...
func (s *State) NewGatewayAuth(id GatewayID, username string, password string) (GatewayAuth, error) {
	if _, exists := s.auths.Gateways[id]; exists {
		return GatewayAuth{}, fmt.Errorf("auth id already exists")
	}
	s.auths.Gateways[id] = &GatewayAuth{
		ID:           id,
		Username:     username,
		passwordHash: hashPassword(password),
	}
	return *s.auths.Gateways[id], nil
}

// DeleteGatewayAuth ...
func (s *State) DeleteGatewayAuth(id GatewayID) {
	delete(s.auths.Gateways, id)
}

// GatewayOTPs ...
func (s *State) GatewayOTPs() []GatewayOTP {
	otps := []GatewayOTP{}
	for _, v := range s.auths.GatewayOTPs {
		otps = append(otps, *v)
	}
	return otps
}

// GatewayOTP ...
func (s *State) GatewayOTP(id GatewayID) (GatewayOTP, error) {
	auth, ok := s.auths.GatewayOTPs[id]
	if !ok {
		return *auth, fmt.Errorf("auth not found")
	}
	return *auth, nil
}

// NewGatewayOTP ...
func (s *State) NewGatewayOTP(id GatewayID, username string, password string, expiratonDuration time.Time) (GatewayOTP, error) {
	if _, exists := s.auths.GatewayOTPs[id]; exists {
		return GatewayOTP{}, fmt.Errorf("auth id already exists")
	}
	s.auths.GatewayOTPs[id] = &GatewayOTP{
		ID:           id,
		Username:     username,
		passwordHash: hashPassword(password),
		Expires:      expiratonDuration,
	}
	return *s.auths.GatewayOTPs[id], nil
}

// DeleteGatewayOTP ...
func (s *State) DeleteGatewayOTP(id GatewayID) {
	delete(s.auths.GatewayOTPs, id)
}

// UserAuths ...
func (s *State) UserAuths() []UserAuth {
	auths := []UserAuth{}
	for _, v := range s.auths.Users {
		auths = append(auths, *v)
	}
	return auths
}

// UserAuth ...
func (s *State) UserAuth(id UserID) (UserAuth, error) {
	auth, ok := s.auths.Users[id]
	if !ok {
		return *auth, fmt.Errorf("auth not found")
	}
	return *auth, nil
}

// NewUserAuth ...
func (s *State) NewUserAuth(id UserID, username string, password string) (UserAuth, error) {
	if _, exists := s.auths.Users[id]; exists {
		return UserAuth{}, fmt.Errorf("auth id already exists")
	}
	s.auths.Users[id] = &UserAuth{
		ID:           id,
		Username:     username,
		passwordHash: hashPassword(password),
	}
	return *s.auths.Users[id], nil
}

// UpdateUserAuthPassword ...
func (s *State) UpdateUserAuthPassword(id UserID, newPassword string) error {
	auth, ok := s.auths.Users[id]
	if !ok {
		return fmt.Errorf("auth not found")
	}
	auth.passwordHash = hashPassword(newPassword)
	return nil
}

// DeleteUserAuth ...
func (s *State) DeleteUserAuth(id UserID) {
	delete(s.auths.Users, id)
}

type config struct {
	GatewayOTPExpirationSeconds int
}

type auths struct {
	Gateways    map[GatewayID]*GatewayAuth
	GatewayOTPs map[GatewayID]*GatewayOTP
	Users       map[UserID]*UserAuth
}

// newAuth is used to avoid nil pointer problems
func newAuths() auths {
	return auths{
		Gateways:    make(map[GatewayID]*GatewayAuth),
		GatewayOTPs: make(map[GatewayID]*GatewayOTP),
	}
}

// TODO:
// - add non-exported functions to create new GatewayAuth/GatewayOTP
// - add exported methods to State that allow for addition of new
//   GatewayAuth/GatewayOTP with validation of relationships. Don't forget
//   locks!

// GatewayAuth ...
type GatewayAuth struct {
	ID           GatewayID
	Username     string
	passwordHash string
}

// GatewayOTP ...
type GatewayOTP struct {
	ID           GatewayID
	Username     string
	passwordHash string
	Expires      time.Time
}

// UserAuth ...
type UserAuth struct {
	ID           UserID
	Username     string
	passwordHash string
}

// currentState ...
type currentState struct {
	Gateways map[GatewayID]*Gateway
	Devices  map[DeviceID]*Device
	Features map[FeatureID]*Feature
}

// newCurrentState is used to avoid nil pointer problems
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

func hashPassword(password string) string {
	// TODO: actually hash
	hash := password
	return hash
}
