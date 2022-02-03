package export

import "time"

type (

	// PlaneLocation is our exported data format. it encodes to JSON
	PlaneLocation struct {
		// This info is populated by the tracker
		New, Removed      bool
		Icao              string
		Lat, Lon, Heading float64
		Velocity          float64
		Altitude          int
		VerticalRate      int
		AltitudeUnits     string
		FlightStatus      string
		OnGround          bool
		Airframe          string
		AirframeType      string
		HasLocation       bool
		HasHeading        bool
		HasVerticalRate   bool
		HasVelocity       bool
		SourceTag         string
		Squawk            string
		Special           string
		TileLocation      string
		TrackedSince      time.Time
		LastMsg           time.Time

		AircraftWidth  *float32 `json:",omitempty"`
		AircraftLength *float32 `json:",omitempty"`

		// Enrichment Plane data
		IcaoCode        *string `json:",omitempty"`
		Registration    *string `json:",omitempty"`
		TypeCode        *string `json:",omitempty"`
		Serial          *string `json:",omitempty"`
		RegisteredOwner *string `json:",omitempty"`
		COFAOwner       *string `json:",omitempty"`
		FlagCode        *string `json:",omitempty"`

		// Enrichment Route Data
		CallSign  *string   `json:",omitempty"`
		Operator  *string   `json:",omitempty"`
		RouteCode *string   `json:",omitempty"`
		Segments  []Segment `json:",omitempty"`
	}

	Segment struct {
		Name     string
		ICAOCode string
	}
)

// Plane here gives us something to look at
func (pl *PlaneLocation) Plane() string {
	if nil != pl.CallSign && "" != *pl.CallSign {
		return *pl.CallSign
	}

	if nil != pl.Registration && "" != *pl.Registration {
		return *pl.Registration
	}

	return "ICAO: " + pl.Icao
}
