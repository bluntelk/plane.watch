package tracker

import "testing"

func TestGlobeIndexSpecialTile_contains(t1 *testing.T) {
	type fields struct {
		North float64
		East  float64
		South float64
		West  float64
	}
	type args struct {
		lat  float64
		long float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"contains centre",
			fields{North: 20, East:  -20, South: -20, West:  20},
			args{lat: 0, long: 0},
			true,
		},
		{
			"world contains Perth",
			fields{North: 90, East:  -180, South: -90, West:  180},
			args{lat: -31.952162, long: 115.943482},
			true,
		},
		{
			"tile contains Perth",
			fields{North: -31, East:  115, South: -32, West:  116},
			args{lat: -31.952162, long: 115.943482},
			true,
		},
		{
			"northern hemisphere does not contain Perth",
			fields{North: 90, East:  -180, South: 0, West:  180},
			args{lat: -31.952162, long: 115.943482},
			false,
		},
		{
			"southern hemisphere does not contain london",
			fields{North: 0, East:  -180, South: -90, West:  180},
			args{lat: 51.5, long: 10},
			false,
		},

		{
			"longitude 0-180 does contain perth",
			fields{North: 90, East:  0, South: -90, West:  180},
			args{lat: -31.952162, long: 115.943482},
			true,
		},
		{
			"longitude -180-0 does not contain perth",
			fields{North: 90, East:  -180, South: -90, West:  0},
			args{lat: -31.952162, long: 115.943482},
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := GlobeIndexSpecialTile{
				North: tt.fields.North,
				East:  tt.fields.East,
				South: tt.fields.South,
				West:  tt.fields.West,
			}
			if got := t.contains(tt.args.lat, tt.args.long); got != tt.want {
				t1.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lookupTile(t *testing.T) {
	type args struct {
		lat float64
		lon float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Perth is found",
			args{lat: -31.952162, lon: 115.943482},
			"tile38",
		},
		{
			"53.253113, 179.723145",
			args{53.253113, 179.723145},
			"tileUnknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lookupTile(tt.args.lat, tt.args.lon); got != tt.want {
				t.Errorf("lookupTile() = %v, want %v", got, tt.want)
			}
		})
	}
}