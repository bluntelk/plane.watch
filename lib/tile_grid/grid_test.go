package tile_grid

import (
	"fmt"
	"testing"
)

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
			fields{North: 20, East: -20, South: -20, West: 20},
			args{lat: 0, long: 0},
			true,
		},
		{
			"world contains Perth",
			fields{North: 90, East: -180, South: -90, West: 180},
			args{lat: -31.952162, long: 115.943482},
			true,
		},
		{
			"tile contains Perth",
			fields{North: -31, East: 115, South: -32, West: 116},
			args{lat: -31.952162, long: 115.943482},
			true,
		},
		{
			"northern hemisphere does not contain Perth",
			fields{North: 90, East: -180, South: 0, West: 180},
			args{lat: -31.952162, long: 115.943482},
			false,
		},
		{
			"southern hemisphere does not contain london",
			fields{North: 0, East: -180, South: -90, West: 180},
			args{lat: 51.5, long: 10},
			false,
		},

		{
			"longitude 0-180 does contain perth",
			fields{North: 90, East: 0, South: -90, West: 180},
			args{lat: -31.952162, long: 115.943482},
			true,
		},
		{
			"longitude -180-0 does not contain perth",
			fields{North: 90, East: -180, South: -90, West: 0},
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

func Test_LookupTile(t *testing.T) {
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
			if got := LookupTile(tt.args.lat, tt.args.lon); got != tt.want {
				t.Errorf("LookupTile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGridLocationNames(t *testing.T) {
	if 0 == len(GridLocationNames()) {
		t.Errorf("Do not have a world")
	}
	if len(worldGrid) != len(GridLocationNames()) {
		t.Errorf("Failed to get the correct number of grid location names")
	}
	for _, name := range GridLocationNames() {
		if "" == name {
			t.Errorf("Got an empty name for tile")
		}
	}
}

func TestTileLookupsSame(t *testing.T) {
	var count, failed int
	for lat := -90.0; lat < 90.0; lat += 1 {
		for lon := -180.0; lon < 180.0; lon += 1 {
			count++
			name := fmt.Sprintf("lookup_%0.2f_%0.2f", lat, lon)
			t.Run(name, func(tt *testing.T) {
				manual := lookupTileManual(lat, lon)
				preCalc := lookupTilePreCalc(lat, lon)
				if manual != preCalc {
					failed++
					tt.Errorf("Lookup Difference. Precalc: %s, manual: %s", preCalc, manual)
				}
			})
		}
	}
	if failed > 0 {
		t.Errorf("%d/%d failed", failed, count)
	}
}

func BenchmarkLookupTileManual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for lat := -90.0; lat < 90.0; lat += 1 {
			for lon := -180.0; lon < 180.0; lon += 1 {
				lookupTileManual(lat, lon)
			}
		}
	}
}

func BenchmarkLookupTilePreCalc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for lat := -90.0; lat < 90.0; lat += 1 {
			for lon := -180.0; lon < 180.0; lon += 1 {
				lookupTilePreCalc(lat, lon)
			}
		}
	}
}
