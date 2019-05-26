package tracker

import (
	"fmt"
	"math"
	"time"
)

// meanings: 0 is even frame, 1 is odd frame
type CprLocation struct {
	evenLat, oddLat, evenLon, oddLon float64

	time0, time1 time.Time

	// working out values
	rlat0, rlat1, airDLat0, airDLat1 float64

	latitudeIndex int32

	oddFrame, evenFrame bool

	nl0, nl1 int32

	globalSurfaceRange float64
}

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



func (cpr *CprLocation) computeLatitudeIndex() {
	cpr.latitudeIndex = int32(math.Floor((((59 * cpr.evenLat) - (60 * cpr.oddLat)) / 131072) + 0.5))
}

func (cpr *CprLocation) computeAirDLatRLat() {
	cpr.airDLat0 = cpr.globalSurfaceRange / 60.0
	cpr.airDLat1 = cpr.globalSurfaceRange / 59.0
	cpr.rlat0 = cpr.airDLat0 * (cprModFunction(cpr.latitudeIndex, 60) + (cpr.evenLat / 131072))
	cpr.rlat1 = cpr.airDLat1 * (cprModFunction(cpr.latitudeIndex, 59) + (cpr.oddLat / 131072))
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
		m := math.Floor((((cpr.evenLon * float64(cpr.nl1-1)) - (cpr.oddLon * float64(cpr.nl1))) / 131072.0) + 0.5)
		//log.Printf("	m = %0.2f", m)

		loc.Longitude = cpr.dlonFunction(cpr.rlat1, 1) * (cprModFunction(int32(m), ni) + (cpr.oddLon / 131072))
		loc.Latitude = cpr.rlat1
		//log.Printf("	rlat = %0.6f, rlon = %0.6f\n", loc.Latitude, loc.Longitude);
	} else {
		// do even decode
		//log.Println("Even Decode")
		ni := cprNFunction(cpr.rlat0, 0)
		//log.Printf("	ni = %d", ni)
		m := math.Floor((((cpr.evenLon * float64(cpr.nl0-1)) - (cpr.oddLon * float64(cpr.nl0))) / 131072) + 0.5)
		//log.Printf("	m = %0.2f", m)
		loc.Longitude = cpr.dlonFunction(cpr.rlat0, 0) * (cprModFunction(int32(m), ni) + cpr.evenLon/131072)
		loc.Latitude = cpr.rlat0
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
		if cpr.oddFrame {
			s = "Have Odd Frame"
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
		cpr.rlat1 -= 90
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
		if cpr.oddFrame {
			s = "Have Odd Frame"
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