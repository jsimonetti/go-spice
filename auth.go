package spice

import (
	"crypto/rand"
	"crypto/rsa"
	"net"

	"crypto/sha1"
)

type Authenticator interface {
	Next(*AuthContext) (accessGranted bool, computeDestination string, err error)
	Method() AuthMethod
	Init() error
}

//go:generate stringer -type=AuthMethod
type AuthMethod uint8

const (
	AuthMethodSpice AuthMethod = 1
	AuthMethodSASL  AuthMethod = 2
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
	tenant          net.Conn
	ticketCrypted   []byte
	ticketUncrypted []byte

	privateKey *rsa.PrivateKey // needed for Spice auth
	otp        string          // previously authenticated ticket
	address    string          // destination compute node
}

func (a *AuthContext) ReadTicket() []byte {
	if len(a.ticketCrypted) != 0 {
		return a.ticketCrypted
	}
	a.ticketCrypted = make([]byte, 128)
	if _, err := a.tenant.Read(a.ticketCrypted); err != nil {
		return nil
	}
	return a.ticketCrypted
}

func (a *AuthContext) PrivateKey() *rsa.PrivateKey {
	key := a.privateKey
	if key == nil {
		return nil
	}
	return key
}

func (a *AuthContext) Password() string {
	crypted := a.ReadTicket()

	key := a.PrivateKey()
	if key == nil {
		return ""
	}

	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha1.New(), rng, key, crypted, []byte{})
	if err != nil {
		return ""
	}

	// do we need to remove this last char??
	a.ticketUncrypted = plaintext[:len(plaintext)-1]

	return string(a.ticketUncrypted)
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
