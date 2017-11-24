package red

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
)

// Various errors which may occur when attempting to marshal or unmarshal
// a SpicePacket to and from its binary form.
var (
	errInvalidPacket  = errors.New("invalid Spice packet")
	errInvalidVersion = errors.New("invalid version")
)

// SpicePacket is the interface used for passing around different kinds of packets.
type SpicePacket interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	validate() error
	finish()
}

func marshalPacket(p SpicePacket) ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func unmarshalPacket(p SpicePacket, b []byte) error {
	buf := bytes.NewReader(b)
	if err := binary.Read(buf, binary.LittleEndian, p); err != nil {
		return err
	}
	return p.validate()
}

var Magic = [4]uint8{0x52, 0x45, 0x44, 0x51}

const (
	VersionMajor uint32 = 2
	VersionMinor uint32 = 2
)

//go:generate stringer -type=ChannelType
type ChannelType uint8

const (
	ChannelMain     ChannelType = 1
	ChannelDisplay  ChannelType = 2
	ChannelInputs   ChannelType = 3
	ChannelCursor   ChannelType = 4
	ChannelPlayback ChannelType = 5
	ChannelRecord   ChannelType = 6
)

//go:generate stringer -type=ErrorCode
type ErrorCode uint32

const (
	ErrorOk                  ErrorCode = 0
	ErrorError               ErrorCode = 1
	ErrorInvalidMagic        ErrorCode = 2
	ErrorInvalidData         ErrorCode = 3
	ErrorVersionMismatch     ErrorCode = 4
	ErrorNeedSecured         ErrorCode = 5
	ErrorNeedUnsecured       ErrorCode = 6
	ErrorPermissionDenied    ErrorCode = 7
	ErrorBadConnectionID     ErrorCode = 8
	ErrorChannelNotAvailable ErrorCode = 9
)

const TicketPubkeyBytes = 162
const ClientTicketBytes = 128

//go:generate stringer -type=Capability
type Capability uint8

const (
	CapabilityAuthSpice Capability = 1
	CapabilityAuthSASL  Capability = 2
)
