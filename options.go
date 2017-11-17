package spice

import (
	"context"
	"net"

	"github.com/sirupsen/logrus"
)

// Option is a functional option handler for Server.
type Option func(*Proxy) error

// SetOption runs a functional option against the server.
func (p *Proxy) SetOption(option Option) error {
	return option(p)
}

func WithLogger(log *logrus.Entry) Option {
	return func(p *Proxy) error {
		p.log = log
		return nil
	}
}

func WithAuthenticator(a Authenticator) Option {
	return func(p *Proxy) error {
		if err := a.Init(); err != nil {
			return err
		}
		p.authenticator[a.Method()] = a
		return nil
	}
}

func WithDialer(dial func(ctx context.Context, network, addr string) (net.Conn, error)) Option {
	return func(p *Proxy) error {
		p.dial = dial
		return nil
	}
}

func defaultDialer() func(context.Context, string, string) (net.Conn, error) {
	dialer := &net.Dialer{}

	return dialer.DialContext
}
