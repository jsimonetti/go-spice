package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jsimonetti/go-spice"
)

func main() {

	log := logrus.New()

	proxy, err := spice.New(spice.WithLogger(log), spice.WithAuthenticator(&AuthSpice{log}))
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	log.Fatal(proxy.ListenAndServe("tcp", "127.0.0.1:5901"))
}

type AuthSpice struct {
	log *logrus.Logger
}

func (a *AuthSpice) Next(ctx *spice.AuthContext) (bool, string, error) {

	ticket := ctx.ReadTicket()
	if ticket == nil {
		err := fmt.Errorf("could not read ticket from client")
		logrus.WithError(err).Error()
		return false, "", err
	}

	key := ctx.PrivateKey()
	if key == nil {
		err := fmt.Errorf("invalid session private key")
		logrus.WithError(err).Error()
		return false, "", err
	}

	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha1.New(), rng, key, ticket, []byte{})
	if err != nil {
		logrus.WithError(err).Error("error decrypting ticket")
		return false, "", fmt.Errorf("error decrypting ticket: %s", err)
	}

	// do we need to remove this last char??
	pass := string(plaintext[:len(plaintext)-1])

	if ctx.OTP() != "" && ctx.OTP() == pass {
		logrus.Debug("OTP found and matches ticket")
		return true, ctx.Address(), nil
	}

	if destination, ok := a.resolveOTPKey(pass); ok {
		logrus.Debug("Ticket validated, compute node at %s", destination)
		ctx.SetOTP(pass)
		ctx.SetAddress(destination)
		return true, ctx.Address(), nil
	}

	logrus.Warn("authentication failed")
	return false, "", nil
}

func (a *AuthSpice) Method() spice.AuthMethod {
	return spice.AuthMethodSpice
}

func (a *AuthSpice) resolveOTPKey(pass string) (string, bool) {
	if pass == "test" {
		logrus.Warn("bogus password check and compute node")
		return "127.0.0.1:5900", true
	}
	return "", false
}

func (a *AuthSpice) Init() error {
	logrus.Debug("AuthSpice initialised")
	return nil
}
