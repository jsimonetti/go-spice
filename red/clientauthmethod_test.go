package red

import (
	"bytes"
	"reflect"
	"testing"
)

func TestClientAuthMethodSelect_UnmarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		cam  ClientAuthMethod
		err  error
	}{
		{
			name: "empty",
			err:  errInvalidPacket,
		},
		{
			name: "short",
			b:    fromHex("00 00 00"),
			err:  errInvalidPacket,
		},
		{
			name: "bad method",
			b:    fromHex("00 00 00 00"),
			err:  errInvalidPacket,
		},
		{
			name: "ok spice",
			cam:  ClientAuthMethod{Method: AuthMethodSpice},
			b:    fromHex("01 00 00 00"),
		},
		{
			name: "ok sasl",
			cam:  ClientAuthMethod{Method: AuthMethodSASL},
			b:    fromHex("02 00 00 00"),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var cam ClientAuthMethod
			err := (&cam).UnmarshalBinary(testCase.b)

			if want, got := testCase.err, err; want != got {
				t.Fatalf("unexpected error:\n- want: %v\n-  got: %v", want, got)
			}
			if err != nil {
				return
			}

			if want, got := testCase.cam, cam; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Message:\n- want: %#v\n-  got: %#v", want, got)
			}
		})
	}
}

func TestClientAuthMethodSelect_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		cam  ClientAuthMethod
		err  error
	}{
		{
			name: "ok spice",
			cam:  ClientAuthMethod{Method: AuthMethodSpice},
			b:    fromHex("01 00 00 00"),
		},
		{
			name: "ok sasl",
			cam:  ClientAuthMethod{Method: AuthMethodSASL},
			b:    fromHex("02 00 00 00"),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			b, err := testCase.cam.MarshalBinary()

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
