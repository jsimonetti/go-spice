package spice

import (
	"fmt"
	"net"

	"context"

	"github.com/sirupsen/logrus"
)

type Proxy struct {
	// WithAuthMethod can be provided to implement custom authentication
	// By default, "auth-less" mode is enabled.
	authenticator map[AuthMethod]Authenticator

	// WithLogger can be used to provide a custom log target.
	// Defaults to stdout.
	log *logrus.Entry

	// WithDialer Optional function for dialing out
	dial func(ctx context.Context, network, addr string) (net.Conn, error)

	// sessionTable
	sessionTable sessionTable
}

func New(options ...Option) (*Proxy, error) {
	proxy := &Proxy{}
	proxy.authenticator = make(map[AuthMethod]Authenticator)

	for _, option := range options {
		if err := proxy.SetOption(option); err != nil {
			return nil, fmt.Errorf("could not set option: %v", err)
		}
	}

	if len(proxy.authenticator) < 1 {
		proxy.authenticator[AuthMethodSpice] = &NOOPAuth{}
	}

	if proxy.log == nil {
		proxy.log = logrus.New().WithField("app", "spiceProxy")
	}

	if proxy.dial == nil {
		proxy.dial = defaultDialer()
	}

	table := sessionTable{}
	table.entries = make(map[SessionID]*sessionEntry)
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
	return nil
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

	flow := NewFlow(tenant, compute)
	if err := flow.Proxy(); err != nil {
		p.log.WithError(err).WithFields(logrus.Fields{"sessionid": handShake.sessionID, "tenant": tenant.RemoteAddr(), "compute": compute.LocalAddr()}).Error("close error")
	}

	p.log.WithFields(logrus.Fields{"sessionid": handShake.sessionID, "tenant": tenant.RemoteAddr(), "compute": compute.LocalAddr()}).Info("connection closed")
	p.sessionTable.Disconnect(handShake.sessionID)

	return nil
}
