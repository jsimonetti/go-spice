package red

type ServerLinkMessage struct {
	// Error codes (i.e., RED_ERROR_?)
	Error ErrorCode

	// PubKey is a 1024 bit RSA public key in X.509 SubjectPublicKeyInfo format
	PubKey [TicketPubkeyBytes]uint8

	// CommonCaps is the number of common client channel capabilities words
	CommonCaps uint32

	// ChannelCaps is the number of specific client channel capabilities words
	ChannelCaps uint32

	// CapsOffset is the location of the start of the capabilities vector given by the
	// bytes offset from the “ size” member (i.e., from the address of the “connection_id”
	// member).
	CapsOffset uint32

	Capabilities1 [4]byte
	Capabilities2 [4]byte
}

// NewServerLinkMessage returns an clientLinkMessage
func NewServerLinkMessage() SpicePacket {
	return &ServerLinkMessage{}
}

// MarshalBinary marshals an ArtPollPacket into a byte slice.
func (p *ServerLinkMessage) MarshalBinary() ([]byte, error) {
	p.finish()
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtPollPacket.
func (p *ServerLinkMessage) UnmarshalBinary(b []byte) error {
	if len(b) < 178 {
		return errInvalidPacket
	}
	full := make([]byte, 186)
	copy(full, b)
	return unmarshalPacket(p, full)
}

// validate is used to validate the Packet.
func (p *ServerLinkMessage) validate() error {
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ServerLinkMessage) finish() {
	p.CapsOffset = 178
}
