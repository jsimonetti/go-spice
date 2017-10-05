package red

type ClientTicket struct {
	Ticket [128]byte
}

// NewClientTicket returns an ClientLinkMessage
func NewClientTicket() SpicePacket {
	return &ClientTicket{}
}

// MarshalBinary marshals an ArtPollPacket into a byte slice.
func (p *ClientTicket) MarshalBinary() ([]byte, error) {
	p.finish()
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtPollPacket.
func (p *ClientTicket) UnmarshalBinary(b []byte) error {
	if len(b) != 128 {
		return errInvalidPacket
	}
	return unmarshalPacket(p, b)
}

// validate is used to validate the Packet.
func (p *ClientTicket) validate() error {
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ClientTicket) finish() {
}
