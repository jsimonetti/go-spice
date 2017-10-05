package red

type ClientLinkMessage struct {
	// Header of the message

	// SessionID In   case   of   a   new   session   (i.e.,   channel   type   is
	// ChannelMain) this field is set to zero, and in response the server will
	// allocate session id and will send it via the RedLinkReply message. In case of all other
	// channel types, this field will be equal to the allocated session id.
	SessionID [4]uint8

	// ChannelType is one of RED_CHANNEL_?
	ChannelType ChannelType

	// ChannelID to connect to. This enables having multiple channels of the same type
	ChannelID uint8

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

// NewClientLinkMessage returns an clientLinkMessage
func NewClientLinkMessage() SpicePacket {
	return &ClientLinkMessage{}
}

// MarshalBinary marshals an ArtPollPacket into a byte slice.
func (p *ClientLinkMessage) MarshalBinary() ([]byte, error) {
	p.finish()
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtPollPacket.
func (p *ClientLinkMessage) UnmarshalBinary(b []byte) error {
	if len(b) < 18 {
		return errInvalidPacket
	}
	full := make([]byte, 26)
	copy(full, b)
	return unmarshalPacket(p, full)
}

// validate is used to validate the Packet.
func (p *ClientLinkMessage) validate() error {
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ClientLinkMessage) finish() {
}
