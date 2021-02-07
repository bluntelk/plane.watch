package mode_s

import "testing"

func Test_inferCommBMessageType(t *testing.T) {
	type args struct {
		mb []byte
	}
	// Logic for the args. If we need to specify bits for detection they are in 0b binary notation. 0's and 0xFF's are junk data
	tests := []struct {
		name    string
		args    args
		want    byte
		want1   byte
		wantErr bool
	}{
		{
			name:    "Correct Length",
			args:    args{mb: []byte{}},
			want:    0,
			want1:   0,
			wantErr: true,
		},
		{
			name:    "Infer BDS 1.0",
			args:    args{mb: []byte{0b0001_0000, 0b1000_0011, 0, 0xFF, 0, 0xFF, 0}},
			want:    1,
			want1:   0,
			wantErr: false,
		},
		{
			name:    "Infer BDS 1.7",
			args:    args{mb: []byte{0b0000_0010, 0xFF, 0xFF, 0b1111_0000, 0b0, 0b0, 0b0}},
			want:    1,
			want1:   7,
			wantErr: false,
		},
		{
			name:    "Infer BDS 2.0",
			args:    args{mb: []byte{0b0010_0000, 0b0100_1100, 0b1001_0000, 0b0111_0010, 0b1100_1011, 0b0100_1000, 0b0010_0000}},
			want:    2,
			want1:   0,
			wantErr: false,
		},
		{
			name:    "Infer BDS 3.0 1",
			args:    args{mb: []byte{0b0011_0000, 0b1111_1110, 0b0011_1100, 0b0000_1000, 0xFF, 0xFF, 0xFF}},
			want:    3,
			want1:   0,
			wantErr: false,
		},
		{
			name:    "Infer BDS 3.0 2",
			args:    args{mb: []byte{0b0011_0000, 0b1111_1110, 0b0100_1100, 0b0000_0100, 0xFF, 0xFF, 0xFF}},
			want:    3,
			want1:   0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := inferCommBMessageType(tt.args.mb)
			if (err != nil) != tt.wantErr {
				t.Errorf("inferCommBMessageType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("inferCommBMessageType() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("inferCommBMessageType() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
