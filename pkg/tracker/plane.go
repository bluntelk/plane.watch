package tracker

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"
)

const (
	max17Bits = 131071
)


type (
PlaneLocation struct {
	Latitude, Longitude  float64
	Altitude             int32
	VerticalRate         int
	AltitudeUnits        string
	Heading, Velocity    float64
	TimeStamp            time.Time
	onGround, hasHeading bool
	hasLatLon            bool
	DistanceTravelled    float64
	DurationTravelled    float64
	TrackFinished        bool
}

 Flight struct {
	Identifier string
	Status     string
	StatusId   byte
}

 Plane struct {
	lastSeen         time.Time
	icaoIdentifier   uint32
	Icao             string
	SquawkIdentity   uint32
	Flight           Flight
	LocationHistory  []PlaneLocation
	Location         PlaneLocation
	cprLocation      CprLocation
	Special          string
	NumUpdates       int
	frameTimes       []time.Time
	RecentFrameCount int
	AirframeCategory string

	rwLock sync.RWMutex
}

 PlaneList map[uint32]Plane

 PlaneIterator func(p Plane)

)
var (
	planeList          PlaneList
	planeAccessMutex   sync.Mutex
	MaxLocationHistory = 10
)

func init() {
	planeList = make(PlaneList, 2000)
}


func (p *Plane) LastSeen() time.Time {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.lastSeen
}

func (p *Plane) SetLastSeen(lastSeen time.Time) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.lastSeen = lastSeen
}

func (p *Plane) IcaoIdentifier() uint32 {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.icaoIdentifier
}

func (p *Plane) SetIcaoIdentifier(icaoIdentifier   uint32) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.icaoIdentifier = icaoIdentifier
}


func GetPlane(ICAO uint32) Plane {
	planeAccessMutex.Lock()
	defer planeAccessMutex.Unlock()

	if plane, ok := planeList[ICAO]; ok {
		return plane
	}
	var p Plane
	p.SetIcaoIdentifier(ICAO)
	p.Icao = fmt.Sprintf("%06X", ICAO)
	p.LocationHistory = make([]PlaneLocation, 0)
	p.ZeroCpr()
	planeList[ICAO] = p

	return p
}

func SetPlane(p Plane, ts time.Time) {
	planeAccessMutex.Lock()
	p.SetLastSeen(ts)
	p.NumUpdates++
	delete(planeList, p.IcaoIdentifier())
	planeList[p.IcaoIdentifier()] = p
	planeAccessMutex.Unlock()
}

func Each(pi PlaneIterator) {
	planeAccessMutex.Lock()
	defer planeAccessMutex.Unlock()

	for _, p := range planeList {
		pi(p)
	}
}

// todo: fix this. it deletes everything
func CleanPlanes() {
	//remove planes that have not been seen for a while
	planeAccessMutex.Lock()
	//
	//// go through the list and remove planes
	//var cutOff, planeCutOff time.Time
	//cutOff = time.Now().Add(5 * time.Minute)
	//for i, plane := range planeList {
	//	planeCutOff = plane.lastSeen.Add(5 * time.Minute)
	//	if planeCutOff.Before(cutOff) {
	//		delete(planeList, i)
	//	}
	//}

	planeAccessMutex.Unlock()
}

func NukePlanes() {
	planeAccessMutex.Lock()
	planeList = make(PlaneList, 2000)
	planeAccessMutex.Unlock()
}

func (p *Plane) String() string {
	var id, alt, position, direction, special, strength string

	white := "\033[0;97m"
	lime := "\033[38;5;118m"
	orange := "\033[38;5;226m"
	blue := "\033[38;5;122m"
	red := "\033[38;5;160m"

	id = fmt.Sprintf("%sPlane (%s%s %-8s%s)", white, lime, p.Icao, p.Flight.Identifier, white)

	if p.Location.onGround {
		position += " is on the ground."
	} else if p.Location.Altitude > 0 {
		alt = fmt.Sprintf(" %s%d%s %s,", orange, p.Location.Altitude, white, p.Location.AltitudeUnits)
	}

	if p.Location.hasLatLon {
		position += fmt.Sprintf(" %s%+03.13f%s, %s%+03.13f%s,", blue, p.Location.Latitude, white, blue, p.Location.Longitude, white)
	}

	if p.Location.hasHeading {
		direction += fmt.Sprintf(" heading %s%0.2f%s, speed %s%0.2f%s knots", orange, p.Location.Heading, white, orange, p.Location.Velocity, white)
	}

	strength = fmt.Sprintf(" %0.2f pps", float64(p.RecentFrameCount)/10.0)

	if "" != p.Special {
		special = " " + red + p.Special + white + ", "
	}

	return id + alt + position + direction + special + strength + "\033[0m"
}

// todo: reimplement as a last seens timestamp? how do we do a count of packets? do we need it?
func (p *Plane) MarkFrameTime(ts time.Time) {
	// cull anything older than 10 seconds
	frameTimes := make([]time.Time, 0)
	cutOff := time.Now().Add(time.Second * -10)
	for _, t := range p.frameTimes {
		if t.After(cutOff) {
			frameTimes = append(frameTimes, t)
		}
	}
	frameTimes = append(frameTimes, ts)
	p.frameTimes = frameTimes
	p.RecentFrameCount = len(p.frameTimes)
}

var PointCounter int

func (p *Plane) AddLatLong(lat, lon float64, ts time.Time) {
	if lat < -95.0 || lat > 95 || lon < -180 || lon > 180 {
		log.Printf("Invalid Coordinate {%0.6f, %0.6f}", lat, lon)
		return
	}

	var distanceTravelled float64
	var durationTravelled float64
	numHistoryItems := len(p.LocationHistory)
	if numHistoryItems > 0 && p.Location.Latitude != 0 && p.Location.Longitude != 0 {
		referenceTime := p.LocationHistory[numHistoryItems-1].TimeStamp
		if !referenceTime.IsZero() {
			durationTravelled = float64(ts.Sub(referenceTime)) / float64(time.Second)
			if 0.0 == durationTravelled {
				durationTravelled = 1
			}
			acceptableMaxDistance := durationTravelled * 343 // mach1 in metres/second seems fast enough...
			if acceptableMaxDistance > 50000 {
				acceptableMaxDistance = 50000
			}

			distanceTravelled = distance(lat, lon, p.Location.Latitude, p.Location.Longitude)

			//log.Printf("%s travelled %0.2fm in %0.2f seconds (%s -> %s)", p.Icao, distanceTravelled, durationTravelled, referenceTime.Format(time.RFC3339Nano), ts.Format(time.RFC3339Nano))

			if distanceTravelled > acceptableMaxDistance {
				log.Printf("The distance (%0.2fm) between {%0.4f,%0.4f} and {%0.4f,%0.4f} is too great for %s to travel in %0.2f seconds. New Track", distanceTravelled, lat, lon, p.Location.Latitude, p.Location.Longitude, p.Icao, durationTravelled)
				p.Location.TrackFinished = true
			}
		}

	}
	PointCounter++

	if MaxLocationHistory > 0 && numHistoryItems >= MaxLocationHistory {
		p.LocationHistory = p.LocationHistory[1:]
	}
	p.LocationHistory = append(p.LocationHistory, p.Location)

	var newLocation PlaneLocation
	newLocation.DistanceTravelled = distanceTravelled
	newLocation.DurationTravelled = durationTravelled
	newLocation.Altitude = p.Location.Altitude
	newLocation.AltitudeUnits = p.Location.AltitudeUnits
	newLocation.hasHeading = p.Location.hasHeading
	newLocation.hasLatLon = true
	newLocation.Heading = p.Location.Heading
	newLocation.Latitude = lat
	newLocation.Longitude = lon
	newLocation.onGround = p.Location.onGround
	newLocation.TimeStamp = ts
	newLocation.Velocity = p.Location.Velocity

	p.Location = newLocation
}

func (p *Plane) ZeroCpr() {
	p.cprLocation.evenLat = 0
	p.cprLocation.evenLon = 0
	p.cprLocation.oddLat = 0
	p.cprLocation.oddLon = 0
	p.cprLocation.rlat0 = 0
	p.cprLocation.rlat1 = 0
	p.cprLocation.time0 = time.Unix(0, 0)
	p.cprLocation.time1 = time.Unix(0, 0)
	p.cprLocation.evenFrame = false
	p.cprLocation.oddFrame = false
}

func (p *Plane) SetCprEvenLocation(lat, lon float64, t time.Time) error {

	// cpr locations are 17 bits long, if we get a value outside of this then we have a problem
	if lat > max17Bits || lat < 0 || lon > max17Bits || lon < 0 {
		return fmt.Errorf("CPR Raw Lat/Lon can be a max of %d, got %0.4f,%0.4f", max17Bits, lat, lon)
	}

	p.cprLocation.evenLat = lat
	p.cprLocation.evenLon = lon
	p.cprLocation.time0 = t
	p.cprLocation.evenFrame = true
	return nil
}

func (p *Plane) SetCprOddLocation(lat, lon float64, t time.Time) error {

	// cpr locations are 17 bits long, if we get a value outside of this then we have a problem
	if lat > max17Bits || lat < 0 || lon > max17Bits || lon < 0 {
		return fmt.Errorf("CPR Raw Lat/Lon can be a max of %d, got %0.4f,%0.4f", max17Bits, lat, lon)
	}

	// only set the odd frame after the even frame is set
	//if !p.cprLocation.evenFrame {
	//	return
	//}

	p.cprLocation.oddLat = lat
	p.cprLocation.oddLon = lon
	p.cprLocation.time1 = t
	p.cprLocation.oddFrame = true
	return nil
}

func (p *Plane) SetAltitude(altitude int32, altitudeUnits string) {
	// set the current altitude
	p.Location.Altitude = altitude
	p.Location.AltitudeUnits = altitudeUnits
}

func (p *Plane) DecodeCpr(ts time.Time) error {

	// attempt to decode the CPR LAT/LON
	var loc PlaneLocation
	var err error

	if p.Location.onGround {
		//loc, err = p.cprLocation.decodeSurface(p.Location.Latitude, p.Location.Longitude)
	} else {
		loc, err = p.cprLocation.decodeGlobalAir()
	}

	if nil != err {
		return err
	}
	p.Location.hasLatLon = true
	p.AddLatLong(loc.Latitude, loc.Longitude, ts)

	//p.ZeroCpr()
	return nil
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
