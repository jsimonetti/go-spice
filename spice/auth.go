package spice

import (
	"context"

	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
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
	ctx context.Context
}

func (a *AuthContext) Ticket() []byte {
	ticket := a.ctx.Value(ContextKeyAuthToken)
	if ticket == nil {
		return nil
	}

	if _, ok := ticket.([]byte); !ok {
		return nil
	}
	return ticket.([]byte)
}

func (a *AuthContext) PrivateKey() *rsa.PrivateKey {
	key := a.ctx.Value(ContextKeyAuthKey)
	if key == nil {
		return nil
	}

	if _, ok := key.(*rsa.PrivateKey); !ok {
		return nil
	}

	return key.(*rsa.PrivateKey)
}

type contextKey string

func (c contextKey) String() string {
	return "spiceAuth context key " + string(c)
}

const (
	ContextKeyAuthToken = contextKey("auth-token")
	ContextKeyAuthKey   = contextKey("auth-key")
)

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

	// crypto/rand.Reader is a good source of entropy for blinding the RSA
	// operation.
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
