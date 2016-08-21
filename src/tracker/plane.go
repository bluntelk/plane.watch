package tracker

import (
	"fmt"
	"math"
	"time"
	"sync"
	//"log"
	"log"
)

const (
	MAX_17_BITS = 131071
)

// meanings: 0 is even frame, 1 is odd frame
type CprLocation struct {
	even_lat, odd_lat, even_lon, odd_lon float64

	time0, time1                         time.Time

	// working out values
	rlat0, rlat1, airDLat0, airDLat1     float64

	latitudeIndex                        int32

	oddFrame, evenFrame                  bool

	nl0, nl1                             int32

	globalSurfaceRange                   float64
}

type PlaneLocation struct {
	Latitude, Longitude  float64
	Altitude             int32
	VerticalRate         int
	AltitudeUnits        string
	Heading, Velocity    float64
	TimeStamp            time.Time
	onGround, hasHeading bool
	hasLatLon            bool
}

type Flight struct {
	Identifier string
	Status     string
	StatusId   byte
}

type Plane struct {
	lastSeen         time.Time
	IcaoIdentifier   uint32
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
}

type PlaneList map[uint32]Plane

type PlaneIterator func(p Plane)

var (
	planeList PlaneList
	planeAccessMutex sync.Mutex
	MaxLocationHistory int = 10
)

var NLTable = map[int32]float64{
	59: 10.47047130,
	58: 14.82817437,
	57: 18.18626357,
	56: 21.02939493,
	55: 23.54504487,
	54: 25.82924707,
	53: 27.93898710,
	52: 29.91135686,
	51: 31.77209708,
	50: 33.53993436,
	49: 35.22899598,
	48: 36.85025108,
	47: 38.41241892,
	46: 39.92256684,
	45: 41.38651832,
	44: 42.80914012,
	43: 44.19454951,
	42: 45.54626723,
	41: 46.86733252,
	40: 48.16039128,
	39: 49.42776439,
	38: 50.67150166,
	37: 51.89342469,
	36: 53.09516153,
	35: 54.27817472,
	34: 55.44378444,
	33: 56.59318756,
	32: 57.72747354,
	31: 58.84763776,
	30: 59.95459277,
	29: 61.04917774,
	28: 62.13216659,
	27: 63.20427479,
	26: 64.26616523,
	25: 65.31845310,
	24: 66.36171008,
	23: 67.39646774,
	22: 68.42322022,
	21: 69.44242631,
	20: 70.45451075,
	19: 71.45986473,
	18: 72.45884545,
	17: 73.45177442,
	16: 74.43893416,
	15: 75.42056257,
	14: 76.39684391,
	13: 77.36789461,
	12: 78.33374083,
	11: 79.29428225,
	10: 80.24923213,
	9:  81.19801349,
	8:  82.13956981,
	7:  83.07199445,
	6:  83.99173563,
	5:  84.89166191,
	4:  85.75541621,
	3:  86.53536998,
	2:  87.00000000,
}

func init() {
	planeList = make(PlaneList, 2000)
}

func GetPlane(ICAO uint32) Plane {
	planeAccessMutex.Lock()

	if plane, ok := planeList[ICAO]; ok {
		planeAccessMutex.Unlock()
		return plane
	}
	var p Plane
	p.IcaoIdentifier = ICAO
	p.Icao = fmt.Sprintf("%06X", ICAO)
	p.LocationHistory = make([]PlaneLocation, 0)
	p.ZeroCpr()
	planeList[ICAO] = p

	planeAccessMutex.Unlock()
	return p
}

func SetPlane(p Plane, ts time.Time) {
	planeAccessMutex.Lock()
	p.lastSeen = ts
	p.NumUpdates++
	delete(planeList, p.IcaoIdentifier)
	planeList[p.IcaoIdentifier] = p
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

	strength = fmt.Sprintf(" %0.2f pps", float64(p.RecentFrameCount) / 10.0)

	if "" != p.Special {
		special = " " + red + p.Special + white + ", "
	}

	return id + alt + position + direction + special + strength + "\033[0m";
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

	duration := float64(ts.Sub(p.Location.TimeStamp) / time.Second)
	println(ts.Format(time.RFC3339Nano))

	if p.Location.Latitude != 0.0 && p.Location.Longitude != 0.0 {
		d := distance(lat, lon, p.Location.Latitude, p.Location.Longitude)
		if d > 5000 { // 5 kilometres
			log.Printf("The distance (%0.2fm) between {%0.4f,%0.4f} and {%0.4f,%0.4f} is too great for %0.2f seconds, ignoring...", d, lat, lon, p.Location.Latitude, p.Location.Longitude, duration)
			return
		}
	}
	PointCounter++

	if MaxLocationHistory > 0 && len(p.LocationHistory) >= MaxLocationHistory {
		p.LocationHistory = p.LocationHistory[1:]
	}
	p.LocationHistory = append(p.LocationHistory, p.Location)

	var newLocation PlaneLocation

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
	p.cprLocation.even_lat = 0
	p.cprLocation.even_lon = 0
	p.cprLocation.odd_lat = 0
	p.cprLocation.odd_lon = 0
	p.cprLocation.rlat0 = 0
	p.cprLocation.rlat1 = 0
	p.cprLocation.time0 = time.Unix(0, 0)
	p.cprLocation.time1 = time.Unix(0, 0)
	p.cprLocation.evenFrame = false
	p.cprLocation.oddFrame = false
}

func (p *Plane) SetCprEvenLocation(lat, lon float64, t time.Time) error {

	// cpr locations are 17 bits long, if we get a value outside of this then we have a problem
	if lat > MAX_17_BITS || lat < 0 || lon > MAX_17_BITS || lon < 0 {
		return fmt.Errorf("CPR Raw Lat/Lon can be a max of %d, got %d,%d", MAX_17_BITS, lat, lon)
	}

	p.cprLocation.even_lat = lat
	p.cprLocation.even_lon = lon
	p.cprLocation.time0 = t
	p.cprLocation.evenFrame = true
	return nil
}

func (p *Plane) SetCprOddLocation(lat, lon float64, t time.Time) error {

	// cpr locations are 17 bits long, if we get a value outside of this then we have a problem
	if lat > MAX_17_BITS || lat < 0 || lon > MAX_17_BITS || lon < 0 {
		return fmt.Errorf("CPR Raw Lat/Lon can be a max of %d, got %d,%d", MAX_17_BITS, lat, lon)
	}

	// only set the odd frame after the even frame is set
	//if !p.cprLocation.evenFrame {
	//	return
	//}

	p.cprLocation.odd_lat = lat
	p.cprLocation.odd_lon = lon
	p.cprLocation.time1 = t
	p.cprLocation.oddFrame = true
	return nil
}

func (p *Plane) SetAltitude(altitude int32, altitude_units string) {
	// set the current altitude
	p.Location.Altitude = altitude
	p.Location.AltitudeUnits = altitude_units
}

func (p *Plane) DecodeCpr(ts time.Time) error {

	// attempt to decode the CPR LAT/LON
	var loc PlaneLocation
	var err error

	if p.Location.onGround {
		loc, err = p.cprLocation.decodeSurface(p.Location.Latitude, p.Location.Longitude)
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

func (cpr *CprLocation) computeLatitudeIndex() {
	cpr.latitudeIndex = int32(math.Floor((((59 * cpr.even_lat) - (60 * cpr.odd_lat)) / 131072) + 0.5))
}

func (cpr *CprLocation) computeAirDLatRLat() {
	cpr.airDLat0 = cpr.globalSurfaceRange / 60.0
	cpr.airDLat1 = cpr.globalSurfaceRange / 59.0
	cpr.rlat0 = cpr.airDLat0 * (cprModFunction(cpr.latitudeIndex, 60) + (cpr.even_lat / 131072))
	cpr.rlat1 = cpr.airDLat1 * (cprModFunction(cpr.latitudeIndex, 59) + (cpr.odd_lat / 131072))
	//log.Printf("j=%d rlat0=%0.6f rlat1=%0.6f", cpr.latitudeIndex, cpr.rlat0, cpr.rlat1)
}

func (cpr *CprLocation) computeLongitudeZone() error {
	cpr.nl0 = getNumLongitudeZone(cpr.rlat0)
	cpr.nl1 = getNumLongitudeZone(cpr.rlat1)

	if cpr.nl0 != cpr.nl1 {
		return fmt.Errorf("Incorrect NL Calculation %d!=%d (for lat/lon %0.13f / %0.13f)", cpr.nl0, cpr.nl1, cpr.rlat0, cpr.rlat1)
	}

	return nil
}

func (cpr *CprLocation) checkFrameTiming() error {
	if cpr.time1.After(cpr.time0.Add(10 * time.Second)) {
		return fmt.Errorf("Unable to decode this CPR Pair. they are too far apart in time (%s, %s)", cpr.time0.Format(time.RFC822Z), cpr.time1.Format(time.RFC822Z))
	}
	return nil
}

func (cpr *CprLocation) computeLatLon() (PlaneLocation, error) {
	var loc PlaneLocation
	if cpr.time1.Before(cpr.time0) {
		//log.Println("Odd Decode")
		// this assumes we are using the odd packet to decode
		/* Compute ni and the longitude index 'm' */
		ni := cprNFunction(cpr.rlat1, 1)
		//log.Printf("	ni = %d", ni)
		m := math.Floor((((cpr.even_lon * float64(cpr.nl1 - 1)) - (cpr.odd_lon * float64(cpr.nl1))) / 131072.0) + 0.5)
		//log.Printf("	m = %0.2f", m)

		loc.Longitude = cpr.dlonFunction(cpr.rlat1, 1) * (cprModFunction(int32(m), ni) + (cpr.odd_lon / 131072))
		loc.Latitude = cpr.rlat1
		//log.Printf("	rlat = %0.6f, rlon = %0.6f\n", loc.Latitude, loc.Longitude);
	} else {
		// do even decode
		//log.Println("Even Decode")
		ni := cprNFunction(cpr.rlat0, 0);
		//log.Printf("	ni = %d", ni)
		m := math.Floor((((cpr.even_lon * float64(cpr.nl0 - 1)) - (cpr.odd_lon * float64(cpr.nl0))) / 131072) + 0.5);
		//log.Printf("	m = %0.2f", m)
		loc.Longitude = cpr.dlonFunction(cpr.rlat0, 0) * (cprModFunction(int32(m), ni) + cpr.even_lon / 131072);
		loc.Latitude = cpr.rlat0;
		//log.Printf("	rlat = %0.6f, rlon = %0.6f\n", loc.Latitude, loc.Longitude);
	}

	if loc.Longitude > 180.0 {
		loc.Longitude -= 180.0
	}
	//log.Printf("post normalise rlat = %0.6f, rlon = %0.6f\n", loc.Latitude, loc.Longitude);

	if loc.Latitude < -90 || loc.Latitude > 90 {
		return PlaneLocation{}, fmt.Errorf("Failed to decode CPR Lat %0.13f is out of range", loc.Latitude)
	}

	return loc, nil
}

func (cpr *CprLocation) decodeSurface(refLat, refLon float64) (PlaneLocation, error) {
	var loc PlaneLocation
	var err error
	cpr.globalSurfaceRange = 90.0

	if 0 == refLat && 0 == refLon {
		return loc, fmt.Errorf("Invalid Reference Location")
	}

	// basic check - make sure we have both frames
	if !(cpr.oddFrame && cpr.evenFrame) {
		var s string
		if (cpr.oddFrame) {
			s = "Have Odd Frame";
		} else {
			s = "Have Even Frame"
		}
		return loc, fmt.Errorf("Need both odd and even frames before decoding, %s", s)
	}

	// Compute the Latitude Index "j"
	cpr.computeLatitudeIndex()
	cpr.computeAirDLatRLat()

	// Pick the quadrant that's closest to the reference location -
	// this is not necessarily the same quadrant that contains the
	// reference location.
	//
	// There are also only two valid quadrants: -90..0 and 0..90;
	// no correct message would try to encoding a latitude in the
	// ranges -180..-90 and 90..180.
	//
	// If the computed latitude is more than 45 degrees north of
	// the reference latitude (using the northern hemisphere
	// solution), then the southern hemisphere solution will be
	// closer to the refernce latitude.
	//
	// e.g. reflat=0, rlat=44, use rlat=44
	//      reflat=0, rlat=46, use rlat=46-90 = -44
	//      reflat=40, rlat=84, use rlat=84
	//      reflat=40, rlat=86, use rlat=86-90 = -4
	//      reflat=-40, rlat=4, use rlat=4
	//      reflat=-40, rlat=6, use rlat=6-90 = -84

	// As a special case, -90, 0 and +90 all encode to zero, so
	// there's a little extra work to do there.

	if cpr.rlat0 == 0 {
		if refLat < -45 {
			cpr.rlat0 = -90
		} else if refLat > 45 {
			cpr.rlat0 = 90
		}
	} else if (cpr.rlat0 - refLat) > 45 {
		cpr.rlat0 -= 90
	}

	if cpr.rlat1 == 0 {
		if refLat < -45 {
			cpr.rlat1 = -90
		} else if refLat > 45 {
			cpr.rlat1 = 90
		}
	} else if (cpr.rlat1 - refLat) > 45 {
		cpr.rlat1 -= 90;
	}

	// Check to see that the latitude is in range: -90 .. +90
	if cpr.rlat0 < -90 || cpr.rlat0 > 90 || cpr.rlat1 < -90 || cpr.rlat1 > 90 {
		return PlaneLocation{}, fmt.Errorf("Failed to decode CPR. Lat out of bounds")
	}

	if err = cpr.computeLongitudeZone(); nil != err {
		return PlaneLocation{}, err
	}

	if err = cpr.checkFrameTiming(); nil != err {
		return PlaneLocation{}, err
	}

	return cpr.computeLatLon()
}

func (cpr *CprLocation) decodeGlobalAir() (PlaneLocation, error) {
	var loc PlaneLocation
	var err error

	// basic check - make sure we have both frames
	if !(cpr.oddFrame && cpr.evenFrame) {
		var s string
		if (cpr.oddFrame) {
			s = "Have Odd Frame";
		} else {
			s = "Have Even Frame"
		}
		return loc, fmt.Errorf("Need both odd and even frames before decoding, %s", s)
	}
	cpr.globalSurfaceRange = 360.0

	// 1. Compute the latitude index (J):
	cpr.computeLatitudeIndex()

	// 2. Compute the values of rlat0 and rlat1:
	cpr.computeAirDLatRLat()

	// Note: Southern hemisphere values are 270° to 360°. Subtract 360°.
	if cpr.rlat0 > 270 {
		cpr.rlat0 = cpr.rlat0 - 360
	}
	if cpr.rlat1 > 270 {
		cpr.rlat1 = cpr.rlat1 - 360
	}

	if err = cpr.computeLongitudeZone(); nil != err {
		return loc, err
	}

	if err = cpr.checkFrameTiming(); nil != err {
		return loc, err
	}

	return cpr.computeLatLon()
}

// NL - The NL function uses the precomputed table from 1090-WP-9-14
func getNumLongitudeZone(lat float64) int32 {
	if lat < 0 {
		lat = math.Abs(lat)
	}

	for i := int32(59); i >= 2; i-- {
		// cannot use range, it does not guarantee order
		if lat < NLTable[i] {
			return i
		}
	}
	return 1
}

func cprNFunction(lat float64, isOdd int32) int32 {
	nl := getNumLongitudeZone(lat) - isOdd
	if nl < 1 {
		nl = 1
	}
	return nl
}

func (cpr *CprLocation) dlonFunction(lat float64, isOdd int32) float64 {
	return cpr.globalSurfaceRange / float64(cprNFunction(lat, isOdd))
}

/* Always positive MOD operation, used for CPR decoding. */
func cprModFunction(a, b int32) float64 {
	res := math.Mod(float64(a), float64(b))
	if res < 0 {
		res += float64(b)
	}
	//return math.Floor(res)
	return res
}

func (pl *PlaneLocation) SetDirection(heading float64, speed int32) {
	pl.Heading = heading
	pl.Velocity = float64(speed)
}



// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
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