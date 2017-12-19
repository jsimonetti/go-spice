package spice

import (
	"context"
	"fmt"
	"net"

	"github.com/jsimonetti/go-spice/red"
	"github.com/sirupsen/logrus"
)

// Proxy is the server object for this spice proxy.
type Proxy struct {
	// WithAuthMethod can be provided to implement custom authentication
	// By default, "auth-less" no-op mode is enabled.
	authenticator map[red.AuthMethod]Authenticator

	// WithLogger can be used to provide a custom logger.
	// Defaults to a logrus implementation.
	log *logrus.Entry

	// WithDialer can be used to provide a custom dialer to reach compute nodes
	// the network is always of type 'tcp' and the computeAddress is the compute node
	// computeAddress that is return by an Authenticator.
	dial func(ctx context.Context, network, addr string) (net.Conn, error)

	// sessionTable holds all the sessions for this proxy
	sessionTable *sessionTable
}

// New returns a new *Proxy with the options applied
func New(options ...Option) (*Proxy, error) {
	proxy := &Proxy{}
	proxy.authenticator = make(map[red.AuthMethod]Authenticator)

	for _, option := range options {
		if err := proxy.SetOption(option); err != nil {
			return nil, fmt.Errorf("could not set option: %v", err)
		}
	}

	if len(proxy.authenticator) < 1 {
		proxy.authenticator[red.AuthMethodSpice] = &NOOPAuth{}
	}

	if proxy.log == nil {
		proxy.log = defaultLogger()
	}

	if proxy.dial == nil {
		proxy.dial = defaultDialer()
	}

	table := &sessionTable{}
	table.entries = make(map[uint32]*sessionEntry)
	proxy.sessionTable = table

	return proxy, nil
}

// ListenAndServe is used to create a listener and serve on it
func (p *Proxy) ListenAndServe(network, addr string) error {
	l, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	return p.Serve(l)
}

// Serve is used to serve connections from a listener
func (p *Proxy) Serve(l net.Listener) error {
	for {
		tenant, err := l.Accept()
		if err != nil {
			return err
		}
		go p.ServeConn(tenant)
	}
}

// ServeConn is used to serve a single connection.
func (p *Proxy) ServeConn(tenant net.Conn) error {
	defer tenant.Close()

	handShake := &tenantHandshake{
		proxy: p,
	}

	var compute net.Conn
	var err error

	for !handShake.Done() {
		if compute, err = handShake.clientLinkStage(tenant); err != nil {
			p.log.WithError(err).Info("handshake failed")
			return err
		}
	}

	p.log.WithFields(logrus.Fields{"sessionid": handShake.sessionID, "tenant": tenant.RemoteAddr(), "compute": compute.LocalAddr()}).Info("connection established")

	flow := newFlow(tenant, compute)
	if err := flow.Pipe(); err != nil {
		p.log.WithError(err).WithFields(logrus.Fields{"sessionid": handShake.sessionID, "tenant": tenant.RemoteAddr(), "compute": compute.LocalAddr()}).Error("close error")
	}

	p.log.WithFields(logrus.Fields{"sessionid": handShake.sessionID, "tenant": tenant.RemoteAddr(), "compute": compute.LocalAddr()}).Info("connection closed")
	p.sessionTable.Disconnect(handShake.sessionID)

	return nil
}
