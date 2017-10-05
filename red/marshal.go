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
	errInvalidPacket = errors.New("invalid Spice packet")
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

type ChannelType uint8

const (
	_ ChannelType = iota
	ChannelMain
	ChannelDisplay
	ChannelInputs
	ChannelCursor
	ChannelPlayback
	ChannelRecord
)

type ErrorCode uint32

const (
	ErrorOk ErrorCode = iota
	ErrorError
	ErrorInvalidMagic
	ErrorInvalidData
	ErrorVersionMismatch
	ErrorNeedSecured
	ErrorNeedUnsecured
	ErrorPermissionDenied
	ErrorBadConnectionID
	ErrorChannelNotAvailable
)

const TicketPubkeyBytes = 162
