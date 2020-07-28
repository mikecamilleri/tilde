package state

// UserInterface ...
type UserInterface struct {
	state *State
}

// NewUserInterface ...
func NewUserInterface(state *State) UserInterface {
	return UserInterface{state: state}
}

// GatewayAuths ...
func (ui *UserInterface) GatewayAuths() []GatewayAuth {
	s := ui.state
	s.RLock()
	defer s.RUnlock()

	gas := []GatewayAuth{}
	for _, gap := range s.auths.Gateways.all() {
		gas = append(gas, *gap)
	}

	return gas
}

// GatewayAuth ...
func (ui *UserInterface) GatewayAuth(id GatewayID) (GatewayAuth, error) {
	s := ui.state
	s.RLock()
	defer s.RUnlock()

	gap, err := s.auths.Gateways.get(id)

	return *gap, err
}

// NewGatewayAuth ...
func (ui *UserInterface) NewGatewayAuth(id GatewayID, username string, password string) (GatewayAuth, error) {
	s := ui.state
	s.Lock()
	defer s.Unlock()

	gap, err := s.auths.Gateways.create(id, username, password)

	return *gap, err
}

// DeleteGatewayAuth ...
func (ui *UserInterface) DeleteGatewayAuth(id GatewayID) {
	s := ui.state
	s.Lock()
	defer s.Unlock()

	s.auths.Gateways.delete(id)
}

// GatewayOTPs ...
func (ui *UserInterface) GatewayOTPs() []GatewayOTP {
	s := ui.state
	s.RLock()
	defer s.RUnlock()

	gos := []GatewayOTP{}
	for _, gop := range s.auths.GatewayOTPs.all() {
		gos = append(gos, *gop)
	}

	return gos
}

// GatewayOTP ...
func (ui *UserInterface) GatewayOTP(id GatewayID) (GatewayOTP, error) {
	s := ui.state
	s.RLock()
	defer s.RUnlock()

	gop, err := s.auths.GatewayOTPs.get(id)

	return *gop, err
}

// NewGatewayOTP ...
func (ui *UserInterface) NewGatewayOTP(id GatewayID, username string, password string) (GatewayOTP, error) {
	s := ui.state
	s.Lock()
	defer s.Unlock()

	gop, err := s.auths.GatewayOTPs.create(id, username, password, s.config.GatewayOTPExpirationSeconds)

	return *gop, err
}

// DeleteGatewayOTP ...
func (ui *UserInterface) DeleteGatewayOTP(id GatewayID) {
	s := ui.state
	s.Lock()
	defer s.Unlock()

	s.auths.GatewayOTPs.delete(id)
}

// UserAuths ...
func (ui *UserInterface) UserAuths() []UserAuth {
	s := ui.state
	s.RLock()
	defer s.RUnlock()

	uas := []UserAuth{}
	for _, uap := range s.auths.Users.all() {
		uas = append(uas, *uap)
	}

	return uas
}

// UserAuth ...
func (ui *UserInterface) UserAuth(id UserID) (UserAuth, error) {
	s := ui.state
	s.Lock()
	defer s.Unlock()

	uap, err := s.auths.Users.get(id)

	return *uap, err
}

// NewUserAuth ...
func (ui *UserInterface) NewUserAuth(id UserID, username string, password string) (UserAuth, error) {
	s := ui.state
	s.Lock()
	defer s.Unlock()

	uap, err := s.auths.Users.create(id, username, password)

	return *uap, err
}

// UpdateUserAuthPassword ...
func (ui *UserInterface) UpdateUserAuthPassword(id UserID, password string) error {
	s := ui.state
	s.Lock()
	defer s.Unlock()

	uap, err := s.auths.Users.get(id)
	if err != nil {
		return err
	}
	uap.setPassword(password)

	return nil
}

// DeleteUserAuth ...
func (ui *UserInterface) DeleteUserAuth(id UserID) {
	s := ui.state
	s.Lock()
	defer s.Unlock()

	s.auths.Users.delete(id)
}
