package red

type RedMiniDataHeader struct {
	// MessageType is type of message
	MessageType uint16

	// Major must be equal to RED_VERSION_MAJOR
	Size uint32
}

// MarshalBinary marshals an ArtPollPacket into a byte slice.
func (p *RedMiniDataHeader) MarshalBinary() ([]byte, error) {
	p.finish()
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into a Packet.
func (p *RedMiniDataHeader) UnmarshalBinary(b []byte) error {
	if len(b) != 6 {
		return errInvalidPacket
	}
	return unmarshalPacket(p, b)
}

func (p *RedMiniDataHeader) validate() error {
	//if !bytes.Equal(p.Magic[:], Magic[:]) {
	//	return errInvalidPacket
	//}
	return nil
}

// finish is used to finish the Packet for sending.
func (p *RedMiniDataHeader) finish() {
}
