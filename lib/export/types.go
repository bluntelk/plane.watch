package export

import "time"

type (
	// EnrichedPlaneLocation is our representation of what the enrichment centre outputs
	EnrichedPlaneLocation struct {
		PlaneLocation  PlaneLocation  `json:"LocationInformation"`
		EnrichmentData EnrichmentData `json:"EnrichmentData"`
	}

	// PlaneLocation is what pw_ingest outputs for its location-updates
	PlaneLocation struct {
		New, Removed      bool
		Icao              string
		Lat, Lon, Heading float64
		Velocity          float64
		Altitude          int
		VerticalRate      int
		AltitudeUnits     string
		FlightNumber      string
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
	}

	EnrichmentData struct {
		Aircraft *Aircraft `json:"aircraft,omitempty"`
	}

	Aircraft struct {
		CofaOwner       *string `json:"cofa_owner"`
		FlagCode        *string `json:"flag_code"`
		IcaoCode        *string `json:"icao_code"`
		RegisteredOwner *string `json:"registered_owner"`
		Registration    *string `json:"registration"`
		Serial          *string `json:"serial"`
		TypeCode        *string `json:"type_code"`
	}
)

func (epl *EnrichedPlaneLocation) Plane() string {
	if "" != epl.PlaneLocation.FlightNumber {
		return epl.PlaneLocation.FlightNumber
	}

	return "ICAO: " + epl.PlaneLocation.Icao
}
