package spice

import (
	"crypto/rsa"
	"net"
)

type Authenticator interface {
	Next(*AuthContext) (accessGranted bool, computeDestination string, err error)
	Method() AuthMethod
	Init() error
}

type AuthMethod uint8

const (
	_ AuthMethod = iota
	AuthMethodSpice
	AuthMethodSASL
)

var _ Authenticator = &NOOPAuth{}

type NOOPAuth struct{}

func (a *NOOPAuth) Next(ctx *AuthContext) (bool, string, error) {
	ctx.ReadTicket()
	return true, "127.0.0.1:5900", nil
}

func (a *NOOPAuth) Method() AuthMethod {
	return AuthMethodSpice
}

func (a *NOOPAuth) Init() error { return nil }

type AuthContext struct {
	tenant     net.Conn
	privateKey *rsa.PrivateKey // needed for Spice auth
	otp        string          // previously authenticated ticket
	address    string          // destination compute node
}

func (a *AuthContext) ReadTicket() []byte {
	ticket := make([]byte, 128)
	if _, err := a.tenant.Read(ticket); err != nil {
		return nil
	}
	return ticket
}

func (a *AuthContext) PrivateKey() *rsa.PrivateKey {
	key := a.privateKey
	if key == nil {
		return nil
	}
	return key
}

func (a *AuthContext) OTP() string {
	return a.otp
}

func (a *AuthContext) SetOTP(otp string) {
	a.otp = otp
}

func (a *AuthContext) Address() string {
	return a.address
}

func (a *AuthContext) SetAddress(address string) {
	a.address = address
}
