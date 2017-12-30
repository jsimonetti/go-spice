package spice

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"log"
	"net"
	"testing"
	"time"
)

func TestAuthSpiceSave(t *testing.T) {
	auth := newAuthSpice(t)

	auth.SaveAddress("123456")
	if auth.LoadAddress() != "123456" {
		t.Errorf("address saved and loaded mismatch")
	}

	auth.SaveToken("123456")
	if auth.LoadToken() != "123456" {
		t.Errorf("tokens saved and loaded mismatch")
	}
}
func TestAuthSpiceToken(t *testing.T) {
	auth := newAuthSpice(t)

	password := "123456"

	// crypto/rand.Reader is a good source of entropy for randomizing the
	// encryption function.
	rng := rand.Reader

	pubkey := auth.privateKey.Public().(*rsa.PublicKey)

	ciphertext, err := rsa.EncryptOAEP(sha1.New(), rng, pubkey, []byte(password), []byte{})
	if err != nil {
		panic(err)
	}
	auth.tenant.Write(ciphertext)

	token, err := auth.Token()
	if err != nil {
		log.Fatalf("unexpected error %#v", err)
	}

	if token != password {
		log.Fatalf("wrong password received")
	}
}

func newAuthSpice(t *testing.T) *authSpice {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	a := &authSpice{
		tenant:     &authConn{},
		privateKey: key,
	}

	return a
}

type authConn struct {
	buf bytes.Buffer
}

func (a *authConn) Read(b []byte) (n int, err error) {
	return a.buf.Read(b)
}

func (a *authConn) Write(b []byte) (n int, err error) {
	return a.buf.Write(b)
}

func (a *authConn) Close() error {
	a.buf.Reset()
	return nil
}

func (authConn) LocalAddr() net.Addr { return nil }

func (authConn) RemoteAddr() net.Addr { return nil }

func (authConn) SetDeadline(t time.Time) error { return nil }

func (authConn) SetReadDeadline(t time.Time) error { return nil }

func (authConn) SetWriteDeadline(t time.Time) error { return nil }
