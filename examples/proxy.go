package main

import (
	"github.com/sirupsen/logrus"

	"github.com/jsimonetti/go-spice"
)

func main() {

	log := logrus.New().WithField("component", "proxy")

	authSpice := &AuthSpice{
		log: log,
	}

	proxy, err := spice.New(spice.WithLogger(log), spice.WithAuthenticator(authSpice))
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	log.Fatal(proxy.ListenAndServe("tcp", "127.0.0.1:5901"))
}

type AuthSpice struct {
	log *logrus.Entry
}

func (a *AuthSpice) Next(ctx *spice.AuthContext) (bool, string, error) {
	pass := ctx.Password()

	if ctx.OTP() != "" && ctx.OTP() == pass {
		a.log.Debug("OTP found and matches password")
		return true, ctx.Address(), nil
	}

	if destination, ok := a.resolveOTPKey(pass); ok {
		a.log.Debug("Ticket validated, compute node at %s", destination)
		ctx.SetOTP(pass)
		ctx.SetAddress(destination)
		return true, ctx.Address(), nil
	}

	a.log.Warn("authentication failed")
	return false, "", nil
}

func (a *AuthSpice) Method() spice.AuthMethod {
	return spice.AuthMethodSpice
}

func (a *AuthSpice) resolveOTPKey(pass string) (string, bool) {
	if pass == "test" {
		a.log.Warn("bogus password check and compute node")
		return "127.0.0.1:5900", true
	}
	return "", false
}

func (a *AuthSpice) Init() error {
	a.log.Debug("AuthSpice initialised")
	return nil
}
