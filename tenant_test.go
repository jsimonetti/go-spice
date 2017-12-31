package spice

import (
	"bytes"
	"crypto"
	"io"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jsimonetti/go-spice/red"
)

func Test_readLinkPacket(t *testing.T) {
	type args struct {
		conn io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				bytes.NewBuffer(make([]byte, 0, 0)),
			},
			wantErr: true,
		},
		{
			name: "no packet",
			args: args{
				bytes.NewBuffer([]byte{
					0x52, 0x45, 0x44, 0x51, 0x02, 0x00, 0x00, 0x00,
					0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
				}),
			},
			wantErr: true,
		},
		{
			name: "1 byte",
			args: args{
				bytes.NewBuffer([]byte{
					0x52, 0x45, 0x44, 0x51, 0x02, 0x00, 0x00, 0x00,
					0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
					0xaa,
				}),
			},
			want: []byte{0xaa},
		},
		{
			name: "1 byte with extra",
			args: args{
				bytes.NewBuffer([]byte{
					0x52, 0x45, 0x44, 0x51, 0x02, 0x00, 0x00, 0x00,
					0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
					0xaa, 0xaa,
				}),
			},
			want: []byte{0xaa},
		},
		{
			name: "8 bytes",
			args: args{
				bytes.NewBuffer([]byte{
					0x52, 0x45, 0x44, 0x51, 0x02, 0x00, 0x00, 0x00,
					0x02, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00,
					0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa,
				}),
			},
			want: []byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readLinkPacket(tt.args.conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("readLinkPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("readLinkPacket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sendServerTicket(t *testing.T) {
	type args struct {
		result red.ErrorCode
	}
	tests := []struct {
		name       string
		args       args
		wantWriter []byte
		wantErr    bool
	}{
		{
			name:       "empty",
			args:       args{},
			wantWriter: []byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			name:       "ok",
			args:       args{red.ErrorOk},
			wantWriter: []byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			name:       "denied",
			args:       args{red.ErrorPermissionDenied},
			wantWriter: []byte{0x07, 0x00, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := bytes.NewBuffer(make([]byte, 0, 0))

			if err := sendServerTicket(tt.args.result, writer); (err != nil) != tt.wantErr {
				t.Errorf("sendServerTicket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.Bytes(); !bytes.Equal(gotWriter, tt.wantWriter) {
				spew.Dump(tt.wantWriter)
				t.Errorf("sendServerTicket() = %+#v, want %+#v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func Test_sendServerLinkPacket(t *testing.T) {
	type args struct {
		key crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		wantWr  string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wr := &bytes.Buffer{}
			if err := sendServerLinkPacket(wr, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("sendServerLinkPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWr := wr.String(); gotWr != tt.wantWr {
				t.Errorf("sendServerLinkPacket() = %v, want %v", gotWr, tt.wantWr)
			}
		})
	}
}

func Test_redPubKey(t *testing.T) {
	type args struct {
		key crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		wantRet red.PubKey
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, err := redPubKey(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("redPubKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("redPubKey() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}
