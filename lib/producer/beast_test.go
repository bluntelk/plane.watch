package producer

import (
	"reflect"
	"testing"
)

var (
	beastModeAc     = []byte{0x1A, 0x31, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	beastModeSShort = []byte{0x1a, 0x32, 0x22, 0x1b, 0x54, 0xf0, 0x81, 0x2b, 0x26, 0x5d, 0x7c, 0x49, 0xf8, 0x28, 0xe9, 0x43}
	beastModeSLong  = []byte{0x1a, 0x33, 0x22, 0x1b, 0x54, 0xac, 0xc2, 0xe9, 0x28, 0x8d, 0x7c, 0x49, 0xf8, 0x58, 0x41, 0xd2, 0x6c, 0xca, 0x39, 0x33, 0xe4, 0x1e, 0xcf}

	beastModeSLongDoubleEsc     = []byte{0x1a, 0x33, 0x22, 0x1b, 0x55, 0xe4, 0x1a, 0x1a, 0xa2, 0x2d, 0x8d, 0x7c, 0x49, 0xf8, 0xe1, 0x1e, 0x2f, 0x00, 0x00, 0x00, 0x00, 0xee, 0xcc, 0x47}
	beastModeSLongDoubleRemoved = []byte{0x1a, 0x33, 0x22, 0x1b, 0x55, 0xe4, 0x1a, 0xa2, 0x2d, 0x8d, 0x7c, 0x49, 0xf8, 0xe1, 0x1e, 0x2f, 0x00, 0x00, 0x00, 0x00, 0xee, 0xcc, 0x47}

	// let's see if we can get a buffer overrun
	beastModeSShortBad = []byte{0xBB, 0x1A, 0x33, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x8D, 0x4D, 0x22, 0x72, 0x99, 0x08, 0x41, 0xB7, 0x90, 0x6C, 0x28, 0x91, 0xA8, 0x1A}
	//                              | ESC | TYPE| MLAT                              | SIG | MODE S LONG
)

func TestScanBeast(t *testing.T) {
	type args struct {
		data  []byte
		atEOF bool
	}
	tests := []struct {
		name        string
		args        args
		wantAdvance int
		wantToken   []byte
		wantErr     bool
	}{
		{
			name:        "Test Not Enough",
			args:        args{data: []byte{0x1a, 0x33, 0x22, 0x1b, 0x54}, atEOF: true},
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     false,
		},
		{
			name:        "Test One Valid Mode AC",
			args:        args{data: beastModeAc, atEOF: true},
			wantAdvance: len(beastModeAc),
			wantToken:   beastModeAc,
			wantErr:     false,
		},
		{
			name:        "Test One Valid Mode S Short",
			args:        args{data: beastModeSShort, atEOF: true},
			wantAdvance: len(beastModeSShort),
			wantToken:   beastModeSShort,
			wantErr:     false,
		},
		{
			name:        "Test One Valid Mode S Long",
			args:        args{data: beastModeSLong, atEOF: true},
			wantAdvance: len(beastModeSLong),
			wantToken:   beastModeSLong,
			wantErr:     false,
		},
		{
			name:        "Test One Valid Mode S Long Double Esc",
			args:        args{data: beastModeSLongDoubleEsc, atEOF: true},
			wantAdvance: len(beastModeSLongDoubleEsc),
			wantToken:   beastModeSLongDoubleRemoved,
			wantErr:     false,
		},
		{
			name:        "Test One Valid Mode S Long Double Esc Buffer Overrun",
			args:        args{data: beastModeSLongDoubleEsc[0:22], atEOF: true},
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     false,
		},
		{
			name:        "Test Two Valid Mode S Short",
			args:        args{data: append(beastModeSShort, beastModeSShort...), atEOF: true},
			wantAdvance: len(beastModeSShort),
			wantToken:   beastModeSShort,
			wantErr:     false,
		},
		{
			name:        "Test Two Valid Mode S Short",
			args:        args{data: append(beastModeSShort[3:], beastModeSShort...), atEOF: true},
			wantAdvance: len(beastModeSShort),
			wantToken:   beastModeSShort,
			wantErr:     false,
		},
		{
			name: "Test Overrun",
			args: args{data: beastModeSShortBad, atEOF: true},
			wantAdvance: 0,
			wantToken: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAdvance, gotToken, err := ScanBeast(tt.args.data, tt.args.atEOF)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanBeast() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAdvance != tt.wantAdvance {
				t.Errorf("ScanBeast() gotAdvance = %v, want %v", gotAdvance, tt.wantAdvance)
			}
			if !reflect.DeepEqual(gotToken, tt.wantToken) {
				t.Errorf("ScanBeast() gotToken (len %d) = %X, want (len %d) %X", len(gotToken), gotToken, len(tt.wantToken), tt.wantToken)
			}
		})
	}
}
