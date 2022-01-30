package mode_s

import (
	"testing"
	"time"
)

func TestFrame_decodeModeSChecksumAddr(t *testing.T) {
	tests := []struct {
		name  string
		frame string
		want  string
		crc   string
	}{
		{
			name:  "df0",
			frame: "*00050319AB8C22;",
			want:  "7C7B5A",
		},
		{
			name:  "df4",
			frame: "*210000992F8C48;",
			want:  "7C7539",
		},
		{
			name:  "df5",
			frame: "28001B1F2181F6;",
			want:  "7C1B28",
		},
		{
			name:  "df16",
			frame: "8061902258822EFC8B9486FDA3BF",
			want:  "7C431F",
		},
		{
			name:  "df20",
			frame: "A000033610020A80F00000270BAA;",
			want:  "7C1666",
		},
		{
			name:  "df21",
			frame: "A80011892058F6B9C38DA09C6D38",
			want:  "7C1BE8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := DecodeString(tt.frame, time.Now())
			if nil != err {
				t.Error(err)
			}
			if nil == f {
				t.Errorf("Failed to decode correctly")
			}
			//t.Logf("%X", f.message)

			if got := f.IcaoStr(); got != tt.want {
				t.Errorf("decodeModeSChecksumAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}
