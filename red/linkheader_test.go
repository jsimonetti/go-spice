package red

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func fromHex(str string) []byte {
	var b []byte
	split := strings.Split(str, " ")
	for _, char := range split {
		c, _ := strconv.ParseUint(char, 16, 8)
		b = append(b, uint8(c))
	}
	return b
}

func TestLinkHeader_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		hdr  LinkHeader
		b    []byte
		err  error
	}{
		{
			name: "empty",
			hdr:  LinkHeader{},
			b:    fromHex("52 45 44 51 02 00 00 00 02 00 00 00 00 00 00 00"),
		},
		{
			name: "len 26",
			hdr:  LinkHeader{Size: 26},
			b:    fromHex("52 45 44 51 02 00 00 00 02 00 00 00 1a 00 00 00"),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			b, err := testCase.hdr.MarshalBinary()

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

func TestLinkHeader_UnmarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		hdr  LinkHeader
		err  error
	}{
		{
			name: "empty",
			err:  errInvalidPacket,
		},
		{
			name: "short",
			b:    fromHex("52 45 44 51 00 00 00 00 02 00 00 00 1a 00 00"),
			err:  errInvalidPacket,
		},
		{
			name: "bad major version",
			err:  errInvalidVersion,
			b:    fromHex("52 45 44 51 00 00 00 00 02 00 00 00 1a 00 00 00"),
		},
		{
			name: "bad minor version",
			err:  errInvalidVersion,
			b:    fromHex("52 45 44 51 02 00 00 00 01 00 00 00 1a 00 00 00"),
		},
		{
			name: "size 26",
			hdr: LinkHeader{
				Magic: Magic,
				Major: VersionMajor,
				Minor: VersionMinor,
				Size:  26,
			},
			b: fromHex("52 45 44 51 02 00 00 00 02 00 00 00 1a 00 00 00"),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var hdr LinkHeader
			err := (&hdr).UnmarshalBinary(testCase.b)

			if want, got := testCase.err, err; want != got {
				t.Fatalf("unexpected error:\n- want: %v\n-  got: %v", want, got)
			}
			if err != nil {
				return
			}

			if want, got := testCase.hdr, hdr; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Message:\n- want: %#v\n-  got: %#v", want, got)
			}
		})
	}
}
