package state

import (
	"fmt"
	"time"
)

// GatewayAuth ...
type GatewayAuth struct {
	ID       GatewayID
	Username string
	password
}

// GatewayOTP ...
type GatewayOTP struct {
	ID       GatewayID
	Expires  time.Time
	Username string
	password
}

// UserAuth ...
type UserAuth struct {
	ID       UserID
	Username string
	password
}

type password struct {
	PasswordHash string `json:"-"`
}

func (pw *password) setPassword(password string) {
	pw.PasswordHash = hashPassword(password)
}

func (pw *password) checkPassword(password string) bool {
	return (pw.PasswordHash == hashPassword(password))
}

func hashPassword(password string) string {
	// TODO: actually hash
	hash := password
	return hash
}

type auths struct {
	Gateways    gatewayAuthsByID
	GatewayOTPs gatewayOTPAuthsByID
	Users       userAuthsByID
}

func newAuths() auths {
	return auths{
		Gateways:    make(gatewayAuthsByID),
		GatewayOTPs: make(gatewayOTPAuthsByID),
		Users:       make(userAuthsByID),
	}
}

type gatewayAuthsByID map[GatewayID]*GatewayAuth

func (m *gatewayAuthsByID) all() []*GatewayAuth {
	gaps := []*GatewayAuth{}
	for _, gap := range *m {
		gaps = append(gaps, gap)
	}
	return gaps
}

func (m *gatewayAuthsByID) get(id GatewayID) (*GatewayAuth, error) {
	gap, ok := (*m)[id]
	if !ok {
		return gap, fmt.Errorf("auth not found")
	}
	return gap, nil
}

func (m *gatewayAuthsByID) create(id GatewayID, username string, password string) (*GatewayAuth, error) {
	if _, ok := (*m)[id]; ok {
		return nil, fmt.Errorf("auth id already exists")
	}
	(*m)[id] = &GatewayAuth{
		ID:       id,
		Username: username,
	}
	(*m)[id].setPassword(password)
	return (*m)[id], nil
}

func (m *gatewayAuthsByID) delete(id GatewayID) {
	delete(*m, id)
}

type gatewayOTPAuthsByID map[GatewayID]*GatewayOTP

func (m *gatewayOTPAuthsByID) all() []*GatewayOTP {
	gops := []*GatewayOTP{}
	for _, gop := range *m {
		gops = append(gops, gop)
	}
	return gops
}

func (m *gatewayOTPAuthsByID) get(id GatewayID) (*GatewayOTP, error) {
	gop, ok := (*m)[id]
	if !ok {
		return nil, fmt.Errorf("auth not found")
	}
	return gop, nil
}

func (m *gatewayOTPAuthsByID) create(id GatewayID, username string, password string, expiratonSeconds int) (*GatewayOTP, error) {
	if _, ok := (*m)[id]; ok {
		return nil, fmt.Errorf("auth id already exists")
	}
	(*m)[id] = &GatewayOTP{
		ID:       id,
		Username: username,
		Expires:  time.Now().Add(time.Duration(expiratonSeconds) * time.Second),
	}
	(*m)[id].setPassword(password)
	return (*m)[id], nil
}

func (m *gatewayOTPAuthsByID) delete(id GatewayID) {
	delete(*m, id)
}

type userAuthsByID map[UserID]*UserAuth

func (m *userAuthsByID) all() []*UserAuth {
	uaps := []*UserAuth{}
	for _, uap := range *m {
		uaps = append(uaps, uap)
	}
	return uaps
}

func (m *userAuthsByID) get(id UserID) (*UserAuth, error) {
	uap, ok := (*m)[id]
	if !ok {
		return nil, fmt.Errorf("auth not found")
	}
	return uap, nil
}

func (m *userAuthsByID) create(id UserID, username string, password string) (*UserAuth, error) {
	if _, ok := (*m)[id]; ok {
		return nil, fmt.Errorf("auth id already exists")
	}
	(*m)[id] = &UserAuth{
		ID:       id,
		Username: username,
	}
	(*m)[id].setPassword(password)
	return (*m)[id], nil
}

func (m *userAuthsByID) delete(id UserID) {
	delete(*m, id)
}
