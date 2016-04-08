package tracker

import (
	"fmt"
	"math"
	"time"
	"sync"
)

const (
	MAX_17_BITS = 131071
)

// meanings: 0 is even frame, 1 is odd frame
type CprLocation struct {
	lat0, lat1, lon0, lon1           float64

	time0, time1                     time.Time

	// working out values
	rlat0, rlat1, airDLat0, airDLat1 float64

	latitudeIndex                    int32

	oddFrame, evenFrame              bool
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
	StatusId   int
}

type Plane struct {
	lastSeen        time.Time
	IcaoIdentifier  uint32
	Icao            string
	SquawkIdentity  uint32
	Flight          Flight
	LocationHistory []PlaneLocation
	Location        PlaneLocation
	cprLocation     CprLocation
	Special         string
}

type PlaneList map[uint32]Plane

var (
	planeList PlaneList
	planeAccessMutex sync.Mutex
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
// NOT WORKING *cry* Instead of fetching an existing entry in the map, it creates a new one
func GetPlane(ICAO uint32) Plane {
	planeAccessMutex.Lock()

	if plane, ok := planeList[ICAO]; ok {
		planeAccessMutex.Unlock()
		return plane
	}
	var p Plane
	p.IcaoIdentifier = ICAO
	p.Icao = fmt.Sprintf("%6X", ICAO)
	p.LocationHistory = make([]PlaneLocation, 0)
	p.ZeroCpr()
	planeList[ICAO] = p

	planeAccessMutex.Unlock()
	return p
}

func SetPlane(p Plane) {
	planeAccessMutex.Lock()
	p.lastSeen = time.Now()
	delete(planeList, p.IcaoIdentifier)
	planeList[p.IcaoIdentifier] = p
	planeAccessMutex.Unlock()
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
func (p *Plane) String() string {
	var id, alt, position, direction string

	id = fmt.Sprintf("\033[0;97mPlane (\033[38;5;118m%s %-8s\033[0;97m) ", p.Icao, p.Flight.Identifier)

	if p.Location.onGround {
		position += " is on the ground. "
	} else if p.Location.Altitude > 0 {
		alt = fmt.Sprintf(" is at %d %s.", p.Location.Altitude, p.Location.AltitudeUnits)
	}

	if p.Location.hasLatLon {
		position += fmt.Sprintf(" at: \033[38;5;122m%+03.13f\033[0;97m, \033[38;5;122m%+03.13f\033[0;97m,", p.Location.Latitude, p.Location.Longitude)
	}

	if p.Location.hasHeading {
		direction += fmt.Sprintf(" Has heading %0.2f and is travelling at %0.2f knots", p.Location.Heading, p.Location.Velocity)
	}

	return id + alt + position + direction + "\033[0m";
}

func (p *Plane) AddLatLong(lat, lon float64) {
	if len(p.LocationHistory) >= 10 {
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
	newLocation.TimeStamp = time.Now()
	newLocation.Velocity = p.Location.Velocity

	p.Location = newLocation
}

func (p *Plane) ZeroCpr() {
	p.cprLocation.lat0 = 0
	p.cprLocation.lon0 = 0
	p.cprLocation.lat1 = 0
	p.cprLocation.lon1 = 0
	p.cprLocation.rlat0 = 0
	p.cprLocation.rlat1 = 0
	p.cprLocation.time0 = time.Unix(0, 0)
	p.cprLocation.time1 = time.Unix(0, 0)
	p.cprLocation.evenFrame = false
	p.cprLocation.oddFrame = false
}

func (p *Plane) SetCprEvenLocation(lat, lon float64, t time.Time) {

	// cpr locations are 17 bits long, if we get a value outside of this then we have a problem
	if lat > MAX_17_BITS || lat < 0 || lon > MAX_17_BITS || lon < 0 {
		return
	}

	p.cprLocation.lat0 = lat
	p.cprLocation.lon0 = lon
	p.cprLocation.time0 = t
	p.cprLocation.evenFrame = true
}

func (p *Plane) SetCprOddLocation(lat, lon float64, t time.Time) {

	// cpr locations are 17 bits long, if we get a value outside of this then we have a problem
	if lat > MAX_17_BITS || lat < 0 || lon > MAX_17_BITS || lon < 0 {
		return
	}

	// only set the odd frame after the even frame is set
	//if !p.cprLocation.evenFrame {
	//	return
	//}

	p.cprLocation.lat1 = lat
	p.cprLocation.lon1 = lon
	p.cprLocation.time1 = t
	p.cprLocation.oddFrame = true
}

func (p *Plane) DecodeCpr(altitude int32, altitude_units string) error {

	// set the current altitude
	p.Location.Altitude = altitude
	p.Location.AltitudeUnits = altitude_units

	// attempt to decode the CPR LAT/LON
	loc, err := p.cprLocation.decode()

	if nil != err {
		return err
	}
	p.Location.hasLatLon = true
	p.AddLatLong(loc.Latitude, loc.Longitude)

	p.ZeroCpr()
	return nil
}

func (cpr *CprLocation) decode() (PlaneLocation, error) {
	var loc PlaneLocation

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

	// 1. Compute the latitude index (J):
	cpr.latitudeIndex = int32((((59 * cpr.lat0) - (60 * cpr.lat1)) / 131072) + 0.5)

	// 2. Compute the values of rlat0 and rlat1:
	cpr.airDLat0 = 360 / 60.0
	cpr.airDLat1 = 360 / 59.0
	cpr.rlat0 = cpr.airDLat0 * (cprModFunction(cpr.latitudeIndex, 60) + (cpr.lat0 / 131072))
	cpr.rlat1 = cpr.airDLat1 * (cprModFunction(cpr.latitudeIndex, 59) + (cpr.lat1 / 131072))

	// Note: Southern hemisphere values are 270° to 360°. Subtract 360°.
	if cpr.rlat0 > 270 {
		cpr.rlat0 = cpr.rlat0 - 360
	}
	if cpr.rlat1 > 270 {
		cpr.rlat1 = cpr.rlat1 - 360
	}

	nl0 := getNumLongitudeZone(cpr.rlat0)
	nl1 := getNumLongitudeZone(cpr.rlat1)

	if nl0 != nl1 {
		return loc, fmt.Errorf("Incorrect NL Calculation %d!=%d (for lat/lon %0.13f / %0.13f)", nl0, nl1, cpr.rlat0, cpr.rlat1)
	}

	if cpr.time1.After(cpr.time0.Add(10 * time.Second)) {
		return loc, fmt.Errorf("Unable to decode this CPR Pair. they are too far apart in time (%s, %s)", cpr.time0.Format(time.RFC822Z), cpr.time1.Format(time.RFC822Z))
	}

	// this assumes we are using the odd packet to decode
	/* Compute ni and the longitude index 'm' */
	ni := cprNFunction(cpr.rlat1, 1)
	m := math.Floor((((cpr.lon0 * float64(nl1 - 1)) - (cpr.lon1 * float64(nl1))) / 131072.0) + 0.5)
	loc.Longitude = cprDlonFunction(cpr.rlat1, 1) * (cprModFunction(int32(m), ni) + (cpr.lon1 / 131072))
	loc.Latitude = cpr.rlat1

	if loc.Longitude > 180.0 {
		loc.Longitude -= 180.0
	}

	return loc, nil
}

// NL - The NL function uses the precomputed table from 1090-WP-9-14
func getNumLongitudeZone(lat float64) int32 {
	if lat < 0 {
		lat = math.Abs(lat)
	}

	for i := int32(59); i > 2; i-- {
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

func cprDlonFunction(lat float64, isOdd int32) float64 {
	return 360.0 / float64(cprNFunction(lat, isOdd))
}

/* Always positive MOD operation, used for CPR decoding. */
func cprModFunction(a, b int32) float64 {
	res := math.Mod(float64(a), float64(b))
	if res < 0 {
		res += float64(b)
	}
	return math.Floor(res)
}

func (pl *PlaneLocation) SetDirection(heading float64, speed int32) {

}
