package spice

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	"io"
	"net"
)

type Authenticator interface {
	Next(AuthContext) (bool, error)
	Method() AuthMethod
	Init() error
}

type AuthMethod uint8

const (
	_ AuthMethod = iota
	AuthMethodSpice
	AuthMethodSASL
)

type NOOPAuth struct{}

func (a *NOOPAuth) Next(ctx AuthContext) (bool, error) {
	return true, nil
}

func (a *NOOPAuth) Method() AuthMethod {
	return AuthMethodSpice
}

func (a *NOOPAuth) Init() error { return nil }

type AuthContext struct {
	tenant     net.Conn
	privateKey *rsa.PrivateKey // needed for Spice auth
}

func (a *AuthContext) Ticket() []byte {
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

func (a *AuthContext) Tenant() io.ReadWriter {
	return a.tenant.(io.ReadWriter)
}

func (a *AuthContext) RemoteAddr() net.Addr {
	return a.tenant.RemoteAddr()
}

func (a *AuthContext) LocalAddr() net.Addr {
	return a.tenant.LocalAddr()
}

type AuthSpice struct{}

func (a *AuthSpice) Next(ctx AuthContext) (bool, error) {

	ticket := ctx.Ticket()
	if ticket == nil {
		return false, fmt.Errorf("unknown ticket")
	}

	key := ctx.PrivateKey()
	if key == nil {
		return false, fmt.Errorf("unknown key")
	}

	rng := rand.Reader

	plaintext, err := rsa.DecryptOAEP(sha1.New(), rng, key, ticket, []byte{})
	if err != nil {
		return false, fmt.Errorf("error in decryption: %s\n", err)
	}

	// do we need to remove this last char??
	pass := string(plaintext[:len(plaintext)-1])
	return a.checkPass(pass), nil
}

func (a *AuthSpice) Method() AuthMethod {
	return AuthMethodSpice
}

func (a *AuthSpice) checkPass(pass string) bool {
	if pass == "test" {
		return true
	}
	return false
}

func (a *AuthSpice) Init() error { return nil }
