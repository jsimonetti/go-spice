package red

import (
	"reflect"
	"testing"
)

func TestClientTicket_UnmarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		ct   ClientTicket
		err  error
	}{
		{
			name: "empty",
			err:  errInvalidPacket,
		},
		{
			name: "short",
			b:    make([]byte, 127),
			err:  errInvalidPacket,
		},
		/*
			{
				name: "ok",
				ct:   ClientTicket{},
				b:    fromHex("61 9a 3d 63 6b e2 b2 b8 4e 94 cb 16 10 b8 ca e0 10 90 a3 01 98 3f b0 fe 8b 3f 6f ca 81 43 15 e8 5b 83 55 4d ae 51 5e ed d7 44 b8 e1 74 25 e6 f7 ba ff 8e fa f9 74 f1 76 a2 9b ea cc 1b 9d b4 3d d7 57 b5 79 11 41 7d fc f6 06 80 0c bb 2e 7f 98 22 46 59 b5 b2 df f8 b7 a3 ad 2a 5b 39 61 24 20 27 d6 17 f8 da a2 e7 f6 ab 1a bc 62 32 e4 ce 5f 91 1b f7 19 0a c1 b9 89 77 7f 01 f1 7d c1 88 54"),
			},
		*/
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var ct ClientTicket
			err := (&ct).UnmarshalBinary(testCase.b)

			if want, got := testCase.err, err; want != got {
				t.Fatalf("unexpected error:\n- want: %v\n-  got: %v", want, got)
			}
			if err != nil {
				return
			}

			if want, got := testCase.ct, ct; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Message:\n- want: %#v\n-  got: %#v", want, got)
			}
		})
	}
}

/*
func TestClientTicket_MarshalBinary(t *testing.T) {
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
*/
