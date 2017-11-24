package red

import "encoding/binary"

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
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b[0:4], uint32(p.Result))
	return b, nil
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtPollPacket.
func (p *ServerTicket) UnmarshalBinary(b []byte) error {
	if len(b) < 4 {
		return errInvalidPacket
	}
	p.Result = ErrorCode(binary.LittleEndian.Uint32(b[0:4]))
	return p.validate()
}

// validate is used to validate the Packet.
func (p *ServerTicket) validate() error {
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ServerTicket) finish() {
}
