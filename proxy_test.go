package spice_test

import (
	"fmt"

	"github.com/jsimonetti/go-spice"
	"github.com/jsimonetti/go-spice/red"
	"github.com/sirupsen/logrus"
)

func ExampleProxy() {
	// create a new logger to be used for the proxy and the authenticator
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	// create a new instance of the sample authenticator
	authSpice := &AuthSpice{
		log: log.WithField("component", "authenticator"),
	}

	// create the proxy using the logger and authenticator
	proxy, err := spice.New(spice.WithLogger(log.WithField("component", "proxy")), spice.WithAuthenticator(authSpice))
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// start listening for tenant connections
	log.Fatal(proxy.ListenAndServe("tcp", "127.0.0.1:5901"))
}

// AuthSpice is an example implementation of a spice Authenticator
type AuthSpice struct {
	log *logrus.Entry

	computeMap map[string]string
}

// Next will check the supplied token and return authorisation information
func (a *AuthSpice) Next(c spice.AuthContext) (bool, string, error) {
	// convert the AuthContext into an AuthSpiceContext, since we do that
	var ctx spice.AuthSpiceContext
	var ok bool
	if ctx, ok = c.(spice.AuthSpiceContext); !ok {
		return false, "", fmt.Errorf("invalid auth method")
	}

	// retrieve the token sent by the tenant
	token, err := ctx.Token()
	if err != nil {
		return false, "", err
	}

	// is the previously saved token is set and matches the token sent by the tenant
	// we return the previously saved compute address
	if ctx.LoadToken() != "" && ctx.LoadToken() == token {
		a.log.Debug("LoadToken found and matches password")
		return true, ctx.LoadAddress(), nil
	}

	// find the compute node for this token
	if destination, ok := a.resolveComputeAddress(token); ok {
		a.log.Debugf("Ticket validated, compute node at %s", destination)
		// save the token and compute address into the context
		// so it can be saved into the session table by the proxy
		ctx.SaveToken(token)
		ctx.SaveAddress(destination)
		return true, ctx.LoadAddress(), nil
	}

	a.log.Warn("authentication failed")
	return false, "", nil
}

// Method returns the Spice auth method
func (a *AuthSpice) Method() red.AuthMethod {
	return red.AuthMethodSpice
}

// resolveComputeAddress is a custom function that checks the token and returns
// a compute node address
func (a *AuthSpice) resolveComputeAddress(token string) (string, bool) {
	// this is just an example, lookup your token somewhere and resolve it to a compute node.
	// When creating your own authentication you should probably use one-time tokens
	// for the tenant authentication. Using a method based on the below sequence of events:
	//
	//  a) Tenant authenticates using token '123e4567:secretpw'
	//  b) The Authenticator looks up the token '123e4567' in a shared store
	//     (kv store or database)
	//  c) The value of token 123e4567 is an encrypted compute node computeAddress.
	//     Attempt to decrypt the computeAddress using 'secretpw'. If this results in a valid
	//     compute node computeAddress, the user is granted access, and de compute destination
	//     is set to the decrypted node computeAddress. In the same transaction, a new
	//     token+secret should be generated, and the old one destroyed

	// lookup in static map
	if compute, ok := a.computeMap[token]; ok {
		a.log.Warn("bogus token check and compute node")
		return compute, true
	}
	return "", false
}

// Init initialises this authenticator
func (a *AuthSpice) Init() error {
	// fill in some compute nodes
	a.computeMap = map[string]string{
		"test":  "127.0.0.1:5900",
		"test2": "127.0.0.1:5902",
	}
	a.log.Debug("AuthSpice initialised")
	return nil
}
