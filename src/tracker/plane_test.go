package tracker

import (
	"fmt"
	"testing"
	"time"
)

func TestFunkyLatLon(t *testing.T) {
	var plane Plane
	var err error
	plane = GetPlane(7777)

	plane.SetCprEvenLocation(92095, 39846, time.Now())
	_, err = plane.cprLocation.decodeAir()
	if nil == err {
		t.Error("We should fail CPR decode with only an even location set")
	}
	plane.ZeroCpr();

	plane.SetCprOddLocation(88385, 125818, time.Now())
	_, err = plane.cprLocation.decodeAir()
	if nil == err {
		t.Error("We should fail CPR decode with only an odd location set")
	}
	plane.ZeroCpr();

	plane = GetPlane(7777)
	plane.SetCprEvenLocation(92095, 39846, time.Now())
	plane.SetCprOddLocation(88385, 125818, time.Now())

	_, err = plane.cprLocation.decodeAir()
	if nil != err {
		t.Error("We should be able to decode with both odd and even CPR locations")
	}
}

func TestGetPlane(t *testing.T) {
	//fmt.Println("TestGetPlane")

	planeListLen := len(planeList)
	var plane Plane
	var err error

	plane = GetPlane(1234)

	if planeListLen == len(planeList) {
		t.Error("Plane List should be longer")
	}

	if 1234 != plane.IcaoIdentifier {
		t.Errorf("Expected planes ICAO identifier to be moo, got %d", plane.IcaoIdentifier)
	}

	SetPlane(plane)

	plane = GetPlane(1234)
	err = plane.SetCprOddLocation(88385, 125818, time.Now())
	if nil != err {
		// there was an error
		t.Errorf("Unexpected error when decoding CPR: %s", err)
	}

	if 88385 != plane.cprLocation.odd_lat {
		t.Errorf("Even Lat not recorded properly. expected 88385, got: %0.2f", plane.cprLocation.odd_lat)
	}

	if 125818 != plane.cprLocation.odd_lon {
		t.Errorf("Even Lon not recorded properly. expected 125818, got: %0.2f", plane.cprLocation.odd_lon)
	}
	SetPlane(plane)

	err = plane.SetCprEvenLocation(92095, 39846, time.Now())
	if nil != err {
		// there was an error
		t.Errorf("Unexpected error when decoding CPR: %s", err)
	}

	if 92095 != plane.cprLocation.even_lat {
		t.Errorf("Even Lat not recorded properly. expected 92095, got: %0.2f", plane.cprLocation.even_lat)
	}

	if 39846 != plane.cprLocation.even_lon {
		t.Errorf("Even Lon not recorded properly. expected 39846, got: %0.2f", plane.cprLocation.even_lon)
	}

	SetPlane(plane)

	plane = GetPlane(1234)
	location, err := plane.cprLocation.decodeAir()

	// ensure the intermediary calculations are correct

	if 1 != plane.cprLocation.latitudeIndex {
		t.Errorf("Incorrect latitude index, expected 1 got %d", plane.cprLocation.latitudeIndex)
	}

	if "10.2157745361328" != fmt.Sprintf("%0.13f", plane.cprLocation.rlat0) {
		t.Errorf("Incorrect RLAT(0) calc, expected 10.2157745361328 - got: %0.13f", plane.cprLocation.rlat0)
	}

	if "10.2162144547802" != fmt.Sprintf("%0.13f", plane.cprLocation.rlat1) {
		t.Errorf("Incorrect RLAT(1) calc, expected 10.2162144547802 - got: %0.13f", plane.cprLocation.rlat1)
	}

	if nil != err {
		// there was an error
		t.Errorf("Unexpected error when decoding CPR: %s", err)
	}

	if "123.889128586342" != fmt.Sprintf("%0.12f", location.Longitude) {
		t.Errorf("Longitude Calculation was incorrect: expected 123.889128586342, got %0.12f", location.Longitude)
	}
	if "10.2162144547802" != fmt.Sprintf("%0.13f", location.Latitude) {
		t.Errorf("Latitude Calculation was incorrect: expected 10.2162144547802, got %0.13f", location.Latitude)
	}

	plane.AddLatLong(location.Latitude, location.Longitude)
	SetPlane(plane)
}

func TestDecodeFailsOnBadData(t *testing.T) {
	plane := GetPlane(1233)
	plane.SetCprEvenLocation(1, 2, time.Now())
	plane.SetCprOddLocation(888888, 888888, time.Now())

	location, err := plane.cprLocation.decodeAir()

	if nil == err {
		t.Errorf("Failed to Fail! we should not be able to decode incomprehensible CPR locations")
	}

	if location.Latitude != 0 {
		t.Errorf("Failed to Fail! we should not be able to decode incomprehensible CPR locations")
	}
}

func TestDecodeFailsOnNoOddLoc(t *testing.T) {
	plane := GetPlane(1235)
	plane.SetCprEvenLocation(92095, 39846, time.Now())

	location, err := plane.cprLocation.decodeAir()

	if nil == err {
		t.Errorf("Failed to Fail! we should not be able to decode when there is no odd CPR location")
	}

	if location.Latitude != 0 {
		t.Errorf("Failed to Fail! we should not be able to decode when there is no odd CPR location")
	}
}
func TestDecodeFailsOnNoEvenLoc(t *testing.T) {
	plane := GetPlane(1236)
	plane.SetCprOddLocation(88385, 125818, time.Now())

	location, err := plane.cprLocation.decodeAir()

	if nil == err {
		t.Errorf("Failed to Fail! we should not be able to decode when there is no even CPR location")
	}

	if location.Latitude != 0 {
		t.Errorf("Failed to Fail! we should not be able to decode when there is no even CPR location")
	}
}
