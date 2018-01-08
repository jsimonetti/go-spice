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
				CommonCapabilities:  []Capability{0x0d},
				ChannelCapabilities: []Capability{0x0f},
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
				CommonCapabilities:  []Capability{0x0d},
				ChannelCapabilities: []Capability{0x0f},
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
