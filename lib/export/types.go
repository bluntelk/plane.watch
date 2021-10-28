package export

import "time"

type (
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
)
