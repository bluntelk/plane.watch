package tile_grid

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"strconv"
)

type (
	GlobalDef struct {
		Refresh                int                     `json:"refresh"`
		History                int                     `json:"history"`
		DbServer               bool                    `json:"dbServer"`
		BinCraft               bool                    `json:"binCraft"`
		GlobeIndexGrid         int                     `json:"globeIndexGrid"`
		GlobeIndexSpecialTiles []GlobeIndexSpecialTile `json:"globeIndexSpecialTiles"`
		Version                string                  `json:"version"`
	}
	GlobeIndexSpecialTile struct {
		North float64 `json:"north"`
		East  float64 `json:"east"`
		South float64 `json:"south"`
		West  float64 `json:"west"`
	}

	GridLocations map[string]GlobeIndexSpecialTile
)

var (
	jsonData  = `{"refresh":1600,"history":1,"dbServer":true,"binCraft":true,"globeIndexGrid":3,"globeIndexSpecialTiles":[{"south":60,"east":0,"north":90,"west":-126},{"south":60,"east":150,"north":90,"west":0},{"south":51,"east":-126,"north":90,"west":150},{"south":9,"east":-126,"north":51,"west":150},{"south":51,"east":-69,"north":60,"west":-126},{"south":45,"east":-114,"north":51,"west":-120},{"south":45,"east":-102,"north":51,"west":-114},{"south":45,"east":-90,"north":51,"west":-102},{"south":45,"east":-75,"north":51,"west":-90},{"south":45,"east":-69,"north":51,"west":-75},{"south":42,"east":18,"north":48,"west":12},{"south":42,"east":24,"north":48,"west":18},{"south":48,"east":24,"north":54,"west":18},{"south":54,"east":24,"north":60,"west":12},{"south":54,"east":12,"north":60,"west":3},{"south":54,"east":3,"north":60,"west":-9},{"south":42,"east":0,"north":48,"west":-9},{"south":42,"east":51,"north":51,"west":24},{"south":51,"east":51,"north":60,"west":24},{"south":30,"east":90,"north":60,"west":51},{"south":30,"east":120,"north":60,"west":90},{"south":30,"east":129,"north":39,"west":120},{"south":30,"east":138,"north":39,"west":129},{"south":30,"east":150,"north":39,"west":138},{"south":39,"east":150,"north":60,"west":120},{"south":9,"east":111,"north":21,"west":90},{"south":21,"east":111,"north":30,"west":90},{"south":9,"east":129,"north":24,"west":111},{"south":24,"east":120,"north":30,"west":111},{"south":24,"east":129,"north":30,"west":120},{"south":9,"east":150,"north":30,"west":129},{"south":9,"east":69,"north":30,"west":51},{"south":9,"east":90,"north":30,"west":69},{"south":-90,"east":51,"north":9,"west":-30},{"south":-90,"east":111,"north":9,"west":51},{"south":-90,"east":160,"north":-18,"west":111},{"south":-18,"east":160,"north":9,"west":111},{"south":-90,"east":-90,"north":-42,"west":160},{"south":-42,"east":-90,"north":9,"west":160},{"south":-9,"east":-42,"north":9,"west":-90},{"south":-90,"east":-63,"north":-9,"west":-90},{"south":-21,"east":-42,"north":-9,"west":-63},{"south":-90,"east":-42,"north":-21,"west":-63},{"south":-90,"east":-30,"north":9,"west":-42},{"south":9,"east":-117,"north":33,"west":-126},{"south":9,"east":-102,"north":30,"west":-117},{"south":9,"east":-90,"north":27,"west":-102},{"south":24,"east":-84,"north":30,"west":-90},{"south":9,"east":-69,"north":18,"west":-90},{"south":18,"east":-69,"north":24,"west":-90},{"south":36,"east":18,"north":42,"west":6},{"south":36,"east":30,"north":42,"west":18},{"south":9,"east":6,"north":39,"west":-9},{"south":9,"east":30,"north":36,"west":6},{"south":9,"east":51,"north":42,"west":30},{"south":24,"east":-69,"north":39,"west":-75},{"south":9,"east":-33,"north":30,"west":-69},{"south":30,"east":-33,"north":60,"west":-69},{"south":9,"east":-9,"north":30,"west":-33},{"south":30,"east":-9,"north":60,"west":-33}],"version":"adsbexchange backend"}`
	worldGrid GridLocations
)

func init() {
	if err := setupWorldGrid(jsonData); nil != err {
		panic(err)
	}
}

func setupWorldGrid(data string) error {
	def := GlobalDef{}
	err := json.Unmarshal([]byte(data), &def)
	if nil != err {
		return err
	}
	worldGrid = make(map[string]GlobeIndexSpecialTile)
	for i, tile := range def.GlobeIndexSpecialTiles {
		worldGrid["tile"+strconv.Itoa(i)] = tile
	}
	return nil
}

func LookupTile(lat, lon float64) string {
	if lat < -95.0 || lat > 95 || lon < -180 || lon > 180 {
		log.Error().Err(fmt.Errorf("cannot lookup invalid coordinates {%0.6f, %0.6f}", lat, lon)).Msg("Using No Tile")
		return ""
	}

	for name, t := range worldGrid {
		if t.contains(lat, lon) {
			return name
		}
	}

	log.Debug().
		Float64("lat", lat).
		Float64("lon", lon).
		Err(fmt.Errorf("could Not Place {%0.6f, %0.6f} in a grid location", lat, lon)).
		Msg("Using No tileUnknown")
	return "tileUnknown"
}

func InGridLocation(lat, lon float64, tileName string) bool {
	if t, ok := worldGrid[tileName]; ok {
		return t.contains(lat, lon)
	}
	return false
}

func GridLocationNames() []string {
	names := make([]string, len(worldGrid))
	i := 0
	for name := range worldGrid {
		names[i] = name
		i++
	}
	return names
}

// contains determines whether the
// * lat is contained between North and South, and
// * lon is contained between East and West
func (t GlobeIndexSpecialTile) contains(lat, lon float64) bool {
	log.Debug().Float64("lat", lat).Float64("north", t.North).Bool("lat <= t.North", lat <= t.North).Send()
	log.Debug().Float64("lat", lat).Float64("south", t.North).Bool("lat >= t.South", lat >= t.South).Send()
	log.Debug().Float64("lon", lat).Float64("east", t.North).Bool("lat >= t.East", lat >= t.East).Send()
	log.Debug().Float64("lon", lat).Float64("west", t.North).Bool("lat <= t.West", lat <= t.West).Send()

	return (lat <= t.North && lat >= t.South) && (lon >= t.East && lon <= t.West)
}

func GetGrid() GridLocations {
	return worldGrid
}
