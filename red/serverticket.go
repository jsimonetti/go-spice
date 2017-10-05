package red

type ServerTicket struct {
	Result ErrorCode
}

// NewServerTicket returns an ClientLinkMessage
func NewServerTicket() SpicePacket {
	return &ServerTicket{}
}

// MarshalBinary marshals an ArtPollPacket into a byte slice.
func (p *ServerTicket) MarshalBinary() ([]byte, error) {
	p.finish()
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtPollPacket.
func (p *ServerTicket) UnmarshalBinary(b []byte) error {
	if len(b) < 4 {
		return errInvalidPacket
	}
	return unmarshalPacket(p, b)
}

// validate is used to validate the Packet.
func (p *ServerTicket) validate() error {
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ServerTicket) finish() {
}
