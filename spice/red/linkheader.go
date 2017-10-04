package red

type RedLinkHeader struct {
	// Magic must be equal to Magic
	Magic [4]uint8

	// Major must be equal to RED_VERSION_MAJOR
	Major uint32

	// Minor must be equal to RED_VERSION_MINOR
	Minor uint32

	// Size in bytes following this field to the end of this message
	Size uint32
}

// MarshalBinary marshals an ArtPollPacket into a byte slice.
func (p *RedLinkHeader) MarshalBinary() ([]byte, error) {
	p.finish()
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into a Packet.
func (p *RedLinkHeader) UnmarshalBinary(b []byte) error {
	if len(b) < 16 {
		return errInvalidPacket
	}
	return unmarshalPacket(p, b)
}

func (p *RedLinkHeader) validate() error {
	//if !bytes.Equal(p.Magic[:], Magic[:]) {
	//	return errInvalidPacket
	//}
	return nil
}

// finish is used to finish the Packet for sending.
func (p *RedLinkHeader) finish() {
	p.Magic = Magic
	p.Major = VersionMajor
	p.Minor = VersionMinor
}
