package tracker

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const (
	max17Bits = 131071
)

type (
	headingInfo []heading
	heading     struct {
		from, to float64
		label    string
	}
	// PlaneLocation stores where we think a plane is currently at. It is am amalgamation of all the tracking info
	// we receive.
	PlaneLocation struct {
		rwlock sync.RWMutex

		latitude, longitude  float64
		altitude             int32
		hasVerticalRate      bool
		verticalRate         int
		altitudeUnits        string
		heading, velocity    float64
		timeStamp            time.Time
		onGround, hasHeading bool
		hasLatLon            bool
		distanceTravelled    float64
		durationTravelled    float64
		TrackFinished        bool
	}

	flight struct {
		identifier string
		status     string
		statusId   byte
	}

	Plane struct {
		tracker          *Tracker
		trackedSince     time.Time
		lastSeen         time.Time
		icaoIdentifier   uint32
		icao             string
		squawk           uint32
		flight           flight
		locationHistory  []*PlaneLocation
		location         *PlaneLocation
		cprLocation      CprLocation
		special          string
		frameTimes       []time.Time
		recentFrameCount int
		msgCount         uint64
		airframeCategory string
		airframeType string

		rwLock sync.RWMutex
	}

	PlaneIterator func(p *Plane) bool

	DistanceTravelled struct {
		metres   float64
		duration float64
	}
)

var (
	MaxLocationHistory = 1000
	headingLookup      = headingInfo{
		{from: 348.75, to: 360, label: "N"},
		{from: 0, to: 11.25, label: "N"},
		{from: 11.25, to: 33.75, label: "NNE"},
		{from: 33.75, to: 56.25, label: "NE"},
		{from: 56.25, to: 78.75, label: "ENE"},
		{from: 78.75, to: 101.25, label: "E"},
		{from: 101.25, to: 123.75, label: "ESE"},
		{from: 123.75, to: 146.25, label: "SE"},
		{from: 146.25, to: 168.75, label: "SSE"},
		{from: 168.75, to: 191.25, label: "S"},
		{from: 191.25, to: 213.75, label: "SSW"},
		{from: 213.75, to: 236.25, label: "SW"},
		{from: 236.25, to: 258.75, label: "WSW"},
		{from: 258.75, to: 281.25, label: "W"},
		{from: 281.25, to: 303.75, label: "WNW"},
		{from: 303.75, to: 326.25, label: "NW"},
		{from: 326.25, to: 348.75, label: "NNW"},
	}
)

func newPlane(icao uint32) *Plane {
	p := &Plane{
		location: &PlaneLocation{},
	}
	p.setIcaoIdentifier(icao)
	p.resetLocationHistory()
	p.zeroCpr()
	p.trackedSince = time.Now()
	return p
}

// TrackedSince tells us when we started tracking this plane (on this run, not historical)
func (p *Plane) TrackedSince() time.Time {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.trackedSince
}

// LastSeen is when we last received a message from this Plane
func (p *Plane) LastSeen() time.Time {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.lastSeen
}

// setLastSeen sets the last seen timestamp
func (p *Plane) setLastSeen(lastSeen time.Time) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.lastSeen = lastSeen
}

// MsgCount is the number of messages we have received from this plane while we have been tracking it
func (p *Plane) MsgCount() uint64 {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.msgCount
}

// incMsgCount increments our message count by 1
func (p *Plane) incMsgCount() {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.msgCount++
}

// IcaoIdentifier returns the ICAO identifier this plane is using
func (p *Plane) IcaoIdentifier() uint32 {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.icaoIdentifier
}

// IcaoIdentifierStr returns a pretty printed ICAO identifier, fit for human consumption
func (p *Plane) IcaoIdentifierStr() string {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.icao
}

// setIcaoIdentifier sets the tracking identifier for this Plane
func (p *Plane) setIcaoIdentifier(icaoIdentifier uint32) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.icaoIdentifier = icaoIdentifier
	p.icao = fmt.Sprintf("%06X", icaoIdentifier)
}

// resetLocationHistory Zeros out the tracking history for this aircraft
func (p *Plane) resetLocationHistory() {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.locationHistory = make([]*PlaneLocation, 0)
}

// setSpecial allows us to set any special status this plane is transmitting
func (p *Plane) setSpecial(status string) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.special = status
}

// Special returns any special status for this aircraft
func (p *Plane) Special() string {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.special
}

// String gives us a nicely printable ANSI escaped string
func (p *Plane) String() string {
	var id, alt, position, direction, special, strength string

	white := "\033[0;97m"
	lime := "\033[38;5;118m"
	orange := "\033[38;5;226m"
	blue := "\033[38;5;122m"
	red := "\033[38;5;160m"

	id = fmt.Sprintf("%sPlane (%s%s %-8s%s)", white, lime, p.IcaoIdentifierStr(), p.FlightNumber(), white)

	if p.OnGround() {
		position += " is on the ground."
	} else if p.Altitude() > 0 {
		alt = fmt.Sprintf(" %s%d%s %s,", orange, p.Altitude(), white, p.AltitudeUnits())
	}

	if p.HasLocation() {
		position += fmt.Sprintf(" %s%+03.13f%s, %s%+03.13f%s,", blue, p.Lat(), white, blue, p.Lon(), white)
	}

	if p.HasHeading() {
		direction += fmt.Sprintf(" heading %s%0.2f%s, speed %s%0.2f%s knots", orange, p.Heading(), white, orange, p.Velocity(), white)
	}

	//strength = fmt.Sprintf(" %0.2f pps", float64(p.recentFrameCount)/10.0)

	if "" != p.special {
		special = " " + red + p.Special() + white + ", "
	}

	return id + alt + position + direction + special + strength + "\033[0m"
}

// setLocationUpdateTime sets the last time the location was updated
func (p *Plane) setLocationUpdateTime(t time.Time) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.location.timeStamp = t
}

// setAltitude puts our plane in the sky
func (p *Plane) setAltitude(altitude int32, altitudeUnits string) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	// set the current altitude

	p.location.altitude = altitude
	p.location.altitudeUnits = altitudeUnits
}

// Altitude is the planes altitude in AltitudeUnits units
func (p *Plane) Altitude() int32 {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	// set the current altitude
	return p.location.altitude
}

// AltitudeUnits how we are measuring altitude (feet / metres)
func (p *Plane) AltitudeUnits() string {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	// set the current altitude
	return p.location.altitudeUnits
}

// setGroundStatus puts our plane on the ground (or not). Use carefully, planes do not like being put on
//the ground suddenly.
func (p *Plane) setGroundStatus(onGround bool) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.location.onGround = onGround
}

// OnGround tells us where the plane thinks it is (In the sky or on the ground)
func (p *Plane) OnGround() bool {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.location.onGround
}

// setFlightStatus sets the flight status of the aircraft, the string is one from mode_s.flightStatusTable
func (p *Plane) setFlightStatus(statusId byte, statusString string) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.flight.statusId = statusId
	p.flight.status = statusString
}

// FlightStatus gives us the flight status of this aircraft
func (p *Plane) FlightStatus() string {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.flight.status
}

// FlightNumber is the planes self identifier for the route it is flying. e.g. QF1, SPTR644
func (p *Plane) FlightNumber() string {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.flight.identifier
}

// setFlightNumber is the flights identifier/number
func (p *Plane) setFlightNumber(flightIdentifier string) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.flight.identifier = flightIdentifier
}

// setSquawkIdentity Sets the planes squawk. A squawk is set by the pilots for various reasons (including flight control)
func (p *Plane) setSquawkIdentity(ident uint32) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.squawk = ident
}

// SquawkIdentity the integer version of the squawk
func (p *Plane) SquawkIdentity() uint32 {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.squawk
}

// SquawkIdentityStr is the string version of SquawkIdentity
func (p *Plane) SquawkIdentityStr() string {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return fmt.Sprint(p.squawk)
}

// setAirFrameCategory is the type of airframe for this aircraft
func (p *Plane) setAirFrameCategory(category string) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.airframeCategory = category
}

func (p *Plane) AirFrame() string {
	return p.airframeCategory
}

// setAirFrameCategory is the type of airframe for this aircraft
func (p *Plane) setAirFrameCategoryType(categoryType string) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.airframeType = categoryType
}

func (p *Plane) AirFrameType() string {
	return p.airframeType
}

// setHeading gives our plane some direction in life
func (p *Plane) setHeading(heading float64) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	// set the current altitude
	p.location.heading = heading
	p.location.hasHeading = true
}

// Heading tells us which way the plane is currently facing
func (p *Plane) Heading() float64 {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	// set the current altitude
	return p.location.heading
}

// HeadingStr gives a nice printable version of the heading, including compass heading
func (p *Plane) HeadingStr() string {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	if !p.location.hasHeading {
		return "?"
	}
	return fmt.Sprintf("%s (%0.2f)", headingLookup.getCompassLabel(p.location.heading), p.location.heading)
}

// HasHeading let's us know if this plane has found it's way in life and knows where it is heading
func (p *Plane) HasHeading() bool {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	// set the current altitude
	return p.location.hasHeading
}

// setVelocity allows us to set the speed the plane is heading
func (p *Plane) setVelocity(velocity float64) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	// set the current altitude
	p.location.velocity = velocity
}

// Velocity is how fast the plane is going in it's Heading
func (p *Plane) Velocity() float64 {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	// set the current altitude
	return p.location.velocity
}
func (p *Plane) VelocityStr() string {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	// set the current altitude
	if p.location.velocity == 0 {
		return "?"
	}
	return fmt.Sprintf("%0.2f knots", p.location.velocity)
}

// DistanceTravelled Tells us how far we have tracked this plane
func (p *Plane) DistanceTravelled() DistanceTravelled {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return DistanceTravelled{
		metres:   p.location.distanceTravelled,
		duration: p.location.durationTravelled,
	}
}

// setVerticalRate shows us how fast the plane is going up and down and uuupp aaannndd doooowwn
func (p *Plane) setVerticalRate(rate int) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.location.hasVerticalRate = true
	p.location.verticalRate = rate
}

// VerticalRate tells us how fast the plane is going up and down
func (p *Plane) VerticalRate() int {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.location.verticalRate
}

// HasVerticalRate tells us if the plane has reported its vertical rate
func (p *Plane) HasVerticalRate() bool {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.location.hasVerticalRate
}

// HasLocation tells us if we have a latitude/longitude decoded
func (p *Plane) HasLocation() bool {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.location.hasLatLon
}

// Lat tells use the planes last reported latitude
func (p *Plane) Lat() float64 {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.location.latitude
}
func (p *Plane) Lon() float64 {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.location.longitude
}


// addLatLong Adds a Lat/Long pair to our location tracking and sets it as the current plane location
func (p *Plane) addLatLong(lat, lon float64, ts time.Time) (warn error) {
	if lat < -95.0 || lat > 95 || lon < -180 || lon > 180 {
		return fmt.Errorf("cannot add invalid coordinates {%0.6f, %0.6f}", lat, lon)
	}
	p.rwLock.Lock()
	defer p.rwLock.Unlock()

	var travelledDistance float64
	var durationTravelled float64
	numHistoryItems := len(p.locationHistory)
	if numHistoryItems > 0 && p.location.latitude != 0 && p.location.longitude != 0 {
		referenceTime := p.locationHistory[numHistoryItems-1].timeStamp
		if !referenceTime.IsZero() {
			durationTravelled = float64(ts.Sub(referenceTime)) / float64(time.Second)
			if 0.0 == durationTravelled {
				durationTravelled = 1
			}
			acceptableMaxDistance := durationTravelled * 343 // mach1 in metres/second seems fast enough...
			if acceptableMaxDistance > 50000 {
				acceptableMaxDistance = 50000
			}

			travelledDistance = distance(lat, lon, p.location.latitude, p.location.longitude)

			//log.Printf("%s travelled %0.2fm in %0.2f seconds (%s -> %s)", p.icaoStr, DistanceTravelled, durationTravelled, referenceTime.Format(time.RFC3339Nano), ts.Format(time.RFC3339Nano))

			if travelledDistance > acceptableMaxDistance {
				warn = fmt.Errorf("the distance (%0.2fm) between {%0.4f,%0.4f} and {%0.4f,%0.4f} is too great for %s to travel in %0.2f seconds. New Track", travelledDistance, lat, lon, p.location.latitude, p.location.longitude, p.icao, durationTravelled)
				p.location.TrackFinished = true
			}
		}
	}

	if MaxLocationHistory > 0 && numHistoryItems >= MaxLocationHistory {
		p.locationHistory = p.locationHistory[1:]
	}
	p.location.latitude = lat
	p.location.longitude = lon
	p.location.hasLatLon = true
	p.locationHistory = append(p.locationHistory, p.location.Copy())
	return
}

// zeroCpr is called once we have successfully decoded our CPR pair
func (p *Plane) zeroCpr() {
	p.cprLocation.zero(true)
}

// setCprEvenLocation sets our Even CPR location for LAT/LON decoding
func (p *Plane) setCprEvenLocation(lat, lon float64, t time.Time) error {
	return p.cprLocation.SetEvenLocation(lat, lon, t)
}

// setCprOddLocation sets our Even CPR location for LAT/LON decoding
func (p *Plane) setCprOddLocation(lat, lon float64, t time.Time) error {
	return p.cprLocation.SetOddLocation(lat, lon, t)
}

// decodeCpr decodes the CPR Even and Odd frames and gets our Plane position
func (p *Plane) decodeCpr(refLat, refLon float64, ts time.Time) error {
	p.cprLocation.refLat = refLat
	p.cprLocation.refLon = refLon
	loc, err := p.cprLocation.decode(p.OnGround())
	if nil != err || loc == nil {
		return err
	}

	err = p.addLatLong(loc.latitude, loc.longitude, ts)
	return err
}

// LocationHistory returns the track history of the Plane
func (p *Plane) LocationHistory() []*PlaneLocation {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.locationHistory
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

// Valid let's us know if we have some data
func (dt *DistanceTravelled) Valid() bool {
	return dt.metres > 0 && dt.duration > 0
}

// Metres returns how far we have gone
func (dt *DistanceTravelled) Metres() float64 {
	return dt.metres
}

// Duration is how long we have been going
func (dt *DistanceTravelled) Duration() float64 {
	return dt.duration
}

// getCompassLabel turns a 0-360 degree compass reading into a nice human readable label N/S/E/W etc
func (hi headingInfo) getCompassLabel(heading float64) string {
	for _, h := range hi {
		if heading >= h.from && heading <= h.to {
			return h.label
		}
	}
	return "?"
}

// Copy let's us duplicate a plane location
func (pl *PlaneLocation) Copy() *PlaneLocation {
	pl.rwlock.RLock()
	defer pl.rwlock.RUnlock()
	return &PlaneLocation{
		latitude:          pl.latitude,
		longitude:         pl.longitude,
		altitude:          pl.altitude,
		hasVerticalRate:   pl.hasVerticalRate,
		verticalRate:      pl.verticalRate,
		altitudeUnits:     pl.altitudeUnits,
		heading:           pl.heading,
		velocity:          pl.velocity,
		timeStamp:         pl.timeStamp,
		onGround:          pl.onGround,
		hasHeading:        pl.hasHeading,
		hasLatLon:         pl.hasLatLon,
		distanceTravelled: pl.distanceTravelled,
		durationTravelled: pl.durationTravelled,
		TrackFinished:     pl.TrackFinished,
	}
}

// Lat returns the Locations current LAT
func (pl *PlaneLocation) Lat() float64 {
	pl.rwlock.RLock()
	defer pl.rwlock.RUnlock()
	return pl.latitude
}

// Lon returns the Locations current LON
func (pl *PlaneLocation) Lon() float64 {
	pl.rwlock.RLock()
	defer pl.rwlock.RUnlock()
	return pl.longitude
}
