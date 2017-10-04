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
	Init()
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

func (a *NOOPAuth) Init() {}

type AuthContext struct {
	ctx context.Context
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

	ticket := ctx.ctx.Value(ContextKeyAuthToken)
	key := ctx.ctx.Value(ContextKeyAuthKey)

	if _, ok := ticket.([]byte); !ok {
		return false, fmt.Errorf("invalid ticket format")
	}
	if _, ok := key.(*rsa.PrivateKey); !ok {
		return false, fmt.Errorf("invalid private key format")
	}

	// crypto/rand.Reader is a good source of entropy for blinding the RSA
	// operation.
	rng := rand.Reader

	plaintext, err := rsa.DecryptOAEP(sha1.New(), rng, key.(*rsa.PrivateKey), ticket.([]byte), []byte{})
	if err != nil {
		return false, fmt.Errorf("Error from decryption: %s\n", err)
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

func (a *AuthSpice) Init() {}
