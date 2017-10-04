package red

type ClientAuthMethodSelect struct {
	Method [4]byte
}

// NewClientAuthMethodSelect returns an ClientLinkMessage
func NewClientAuthMethodSelect() SpicePacket {
	return &ClientAuthMethodSelect{}
}

// MarshalBinary marshals an ArtPollPacket into a byte slice.
func (p *ClientAuthMethodSelect) MarshalBinary() ([]byte, error) {
	p.finish()
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtPollPacket.
func (p *ClientAuthMethodSelect) UnmarshalBinary(b []byte) error {
	if len(b) < 4 {
		return errInvalidPacket
	}
	return unmarshalPacket(p, b)
}

// validate is used to validate the Packet.
func (p *ClientAuthMethodSelect) validate() error {
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ClientAuthMethodSelect) finish() {
}
