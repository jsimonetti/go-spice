package red

import (
	"bytes"
	"reflect"
	"testing"
)

func TestClientLinkMessage_UnmarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		clm  ClientLinkMessage
		err  error
	}{
		{
			name: "no caps",
			clm: ClientLinkMessage{
				SessionID:           0,
				ChannelType:         1,
				ChannelID:           0,
				CommonCaps:          0,
				ChannelCaps:         0,
				CapsOffset:          18,
				CommonCapabilities:  nil,
				ChannelCapabilities: nil,
			},
			b: fromHex("00 00 00 00 01 00 00 00 00 00 00 00 00 00 12 00 00 00"),
		},
		{
			name: "ok",
			clm: ClientLinkMessage{
				SessionID:           0,
				ChannelType:         1,
				ChannelID:           0,
				CommonCaps:          1,
				ChannelCaps:         1,
				CapsOffset:          18,
				CommonCapabilities:  []uint32{0x0d},
				ChannelCapabilities: []uint32{0x0f},
			},
			b: fromHex("00 00 00 00 01 00 01 00 00 00 01 00 00 00 12 00 00 00 0d 00 00 00 0f 00 00 00"),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var clm ClientLinkMessage
			err := (&clm).UnmarshalBinary(testCase.b)

			if want, got := testCase.err, err; want != got {
				t.Fatalf("unexpected error:\n- want: %v\n-  got: %v", want, got)
			}
			if err != nil {
				return
			}

			if want, got := testCase.clm, clm; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Message:\n- want: %#v\n-  got: %#v", want, got)
			}
		})
	}
}

func TestClientLinkMessage_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		clm  ClientLinkMessage
		err  error
	}{
		{
			name: "empty",
			err:  errInvalidPacket,
		},
		{
			name: "short",
			b:    fromHex("00 00 00 00 01 00 00 00 00 00 00 00 00 00 12 00 00"),
			err:  errInvalidPacket,
		},
		{
			name: "uneven caps",
			b:    fromHex("00 00 00 00 01 00 00 00 00 00 10 00 00 00 12 00 00 00"),
			err:  errInvalidPacket,
		},

		{
			name: "no caps",
			clm: ClientLinkMessage{
				SessionID:           0,
				ChannelType:         1,
				ChannelID:           0,
				CommonCaps:          0,
				ChannelCaps:         0,
				CapsOffset:          18,
				CommonCapabilities:  nil,
				ChannelCapabilities: nil,
			},
			b: fromHex("00 00 00 00 01 00 00 00 00 00 00 00 00 00 12 00 00 00"),
		},
		{
			name: "wrong caps",
			b:    fromHex("00 00 00 00 30 81 9f 30 0d 06 09 2a 86 48 86 f7 0d 01 01 01 05 00 03 81 8d 00 30 81 89 02 81 81 00 bb 49 20 f1 9e 70 a3 07 32 ca a1 63 ce 8d 05 26 82 73 3a 74 59 9d cc c3 83 9c c8 59 60 7e 15 5b 62 8d 53 02 aa f4 81 bf e6 b5 bc 17 88 10 4c d6 dc 6c 83 b9 c2 05 4e ed 89 99 a7 a3 fd 2d 05 d3 0d 60 b3 de 6d 16 3c 9e c8 8c 33 38 b8 3d 39 c1 23 d7 c3 ae e0 59 b6 1a b1 87 d5 b5 30 dc 2b 04 c7 92 6d 92 c4 be bf 21 ae 8a 69 ff 53 1c 41 ff a7 1d 32 8d bb 86 aa c2 50 c4 da 53 f9 24 b0 99 02 03 01 00 01 01 00 00 00 01 00 00 00 b2 00 00 00"),
			err:  errInvalidPacket,
		},
		{
			name: "ok",
			clm: ClientLinkMessage{
				SessionID:           0,
				ChannelType:         1,
				ChannelID:           0,
				CommonCaps:          1,
				ChannelCaps:         1,
				CapsOffset:          18,
				CommonCapabilities:  []uint32{0x0d},
				ChannelCapabilities: []uint32{0x0f},
			},
			b: fromHex("00 00 00 00 01 00 01 00 00 00 01 00 00 00 12 00 00 00 0d 00 00 00 0f 00 00 00"),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			b, err := testCase.clm.MarshalBinary()

			if want, got := testCase.err, err; want != got {
				t.Fatalf("unexpected error:\n- want: %v\n-  got: %v", want, got)
			}
			if err != nil {
				return
			}

			if want, got := testCase.b, b; !bytes.Equal(want, got) {
				t.Fatalf("unexpected Message bytes:\n- want: [%# x]\n-  got: [%# x]", want, got)
			}
		})
	}
}
