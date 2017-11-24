package red

import "encoding/binary"

type ClientAuthMethod struct {
	Method AuthMethod
}

// NewClientAuthMethodSelect returns an ClientLinkMessage
func NewClientAuthMethodSelect() SpicePacket {
	return &ClientAuthMethod{}
}

// MarshalBinary marshals an ArtPollPacket into a byte slice.
func (p *ClientAuthMethod) MarshalBinary() ([]byte, error) {
	p.finish()
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b[0:4], uint32(p.Method))
	return b, nil
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtPollPacket.
func (p *ClientAuthMethod) UnmarshalBinary(b []byte) error {
	if len(b) < 4 {
		return errInvalidPacket
	}
	p.Method = AuthMethod(binary.LittleEndian.Uint32(b[0:4]))
	return p.validate()
}

// validate is used to validate the Packet.
func (p *ClientAuthMethod) validate() error {
	if p.Method != AuthMethodSpice && p.Method != AuthMethodSASL {
		return errInvalidPacket
	}
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ClientAuthMethod) finish() {
}

//go:generate stringer -type=AuthMethod
type AuthMethod uint32

const (
	AuthMethodSpice AuthMethod = 1
	AuthMethodSASL  AuthMethod = 2
)
