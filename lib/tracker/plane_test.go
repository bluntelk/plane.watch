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
	_, err = plane.cprLocation.decodeGlobalAir()
	if nil == err {
		t.Error("We should fail CPR decode with only an even location set")
	}
	plane.ZeroCpr();

	plane.SetCprOddLocation(88385, 125818, time.Now())
	_, err = plane.cprLocation.decodeGlobalAir()
	if nil == err {
		t.Error("We should fail CPR decode with only an odd location set")
	}
	plane.ZeroCpr();

	plane = GetPlane(7777)
	plane.SetCprEvenLocation(92095, 39846, time.Now())
	plane.SetCprOddLocation(88385, 125818, time.Now())

	_, err = plane.cprLocation.decodeGlobalAir()
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

	if 1234 != plane.IcaoIdentifier() {
		t.Errorf("Expected planes ICAO identifier to be moo, got %d", plane.IcaoIdentifier())
	}

	SetPlane(plane, time.Now())

	plane = GetPlane(1234)
	err = plane.SetCprOddLocation(88385, 125818, time.Now())
	if nil != err {
		// there was an error
		t.Errorf("Unexpected error when decoding CPR: %s", err)
	}

	if 88385 != plane.cprLocation.oddLat {
		t.Errorf("Even Lat not recorded properly. expected 88385, got: %0.2f", plane.cprLocation.oddLat)
	}

	if 125818 != plane.cprLocation.oddLon {
		t.Errorf("Even Lon not recorded properly. expected 125818, got: %0.2f", plane.cprLocation.oddLon)
	}
	SetPlane(plane, time.Now())

	err = plane.SetCprEvenLocation(92095, 39846, time.Now())
	if nil != err {
		// there was an error
		t.Errorf("Unexpected error when decoding CPR: %s", err)
	}

	if 92095 != plane.cprLocation.evenLat {
		t.Errorf("Even Lat not recorded properly. expected 92095, got: %0.2f", plane.cprLocation.evenLat)
	}

	if 39846 != plane.cprLocation.evenLon {
		t.Errorf("Even Lon not recorded properly. expected 39846, got: %0.2f", plane.cprLocation.evenLon)
	}

	SetPlane(plane, time.Now())

	plane = GetPlane(1234)
	location, err := plane.cprLocation.decodeGlobalAir()

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

	plane.AddLatLong(location.Latitude, location.Longitude, time.Now())
	SetPlane(plane, time.Now())
}

func TestDecodeFailsOnBadData(t *testing.T) {
	plane := GetPlane(1233)
	plane.SetCprEvenLocation(1, 2, time.Now())
	plane.SetCprOddLocation(888888, 888888, time.Now())

	location, err := plane.cprLocation.decodeGlobalAir()

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

	location, err := plane.cprLocation.decodeGlobalAir()

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

	location, err := plane.cprLocation.decodeGlobalAir()

	if nil == err {
		t.Errorf("Failed to Fail! we should not be able to decode when there is no even CPR location")
	}

	if location.Latitude != 0 {
		t.Errorf("Failed to Fail! we should not be able to decode when there is no even CPR location")
	}
}

func TestCprDecodeSurfacePosition(t *testing.T) {

	type surfaceTestTable struct {
		refLat, refLon           float64
		even_cprLat, even_cprLon float64
		odd_cprLat, odd_cprLon   float64

		evenErrCount             int
		even_rLat, even_rLon     float64
		oddErrCount              int
		odd_rLat, odd_rLon       float64
	}

	// yanked from mutability's dump1090 cprtests.c
	testData := []surfaceTestTable{
		{52.00, -180.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601 - 180.0, 0, 52.209976, 0.176507 - 180.0 },
		{52.00, -140.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601 - 180.0, 0, 52.209976, 0.176507 - 180.0 },
		{52.00, -130.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601 - 90.0, 0, 52.209976, 0.176507 - 90.0 },
		{52.00, -50.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601 - 90.0, 0, 52.209976, 0.176507 - 90.0 },
		{52.00, -40.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601, 0, 52.209976, 0.176507 },
		{52.00, -10.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601, 0, 52.209976, 0.176507 },
		{52.00, 0.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601, 0, 52.209976, 0.176507 },
		{52.00, 10.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601, 0, 52.209976, 0.176507 },
		{52.00, 40.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601, 0, 52.209976, 0.176507 },
		{52.00, 50.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601 + 90.0, 0, 52.209976, 0.176507 + 90.0 },
		{52.00, 130.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601 + 90.0, 0, 52.209976, 0.176507 + 90.0 },
		{52.00, 140.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601 - 180.0, 0, 52.209976, 0.176507 - 180.0 },
		{52.00, 180.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601 - 180.0, 0, 52.209976, 0.176507 - 180.0 },

		// latitude quadrants (but only 2). The decoded longitude also changes because the cell size changes with latitude
		{90.00, 0.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601, 0, 52.209976, 0.176507 },
		{52.00, 0.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601, 0, 52.209976, 0.176507 },
		{8.00, 0.00, 105730, 9259, 29693, 8997, 0, 52.209984, 0.176601, 0, 52.209976, 0.176507 },
		{7.00, 0.00, 105730, 9259, 29693, 8997, 0, 52.209984 - 90.0, 0.135269, 0, 52.209976 - 90.0, 0.134299 },
		{-52.00, 0.00, 105730, 9259, 29693, 8997, 0, 52.209984 - 90.0, 0.135269, 0, 52.209976 - 90.0, 0.134299 },
		{-90.00, 0.00, 105730, 9259, 29693, 8997, 0, 52.209984 - 90.0, 0.135269, 0, 52.209976 - 90.0, 0.134299 },

		// poles/equator cases
		{-46.00, -180.00, 0, 0, 0, 0, 0, -90.0, -180.000000, 0, -90.0, -180.0 }, // south pole
		{-44.00, -180.00, 0, 0, 0, 0, 0, 0.0, -180.000000, 0, 0.0, -180.0 }, // equator
		{44.00, -180.00, 0, 0, 0, 0, 0, 0.0, -180.000000, 0, 0.0, -180.0 }, // equator
		{46.00, -180.00, 0, 0, 0, 0, 0, 90.0, -180.000000, 0, 90.0, -180.0 }, // north pole
	}
	var plane Plane
	var loc PlaneLocation
	var err error
	var expectedLat, expectedLon, actualLat, actualLon string

	for i, test := range testData {
		NukePlanes()
		plane = GetPlane(99887)
		plane.SetCprEvenLocation(test.even_cprLat, test.even_cprLon, time.Now())
		plane.SetCprOddLocation(test.odd_cprLat, test.odd_cprLon, time.Now())
		loc, err = plane.cprLocation.decodeSurface(test.refLat, test.refLon)

		if nil != err && test.evenErrCount == 0 {
			t.Error(err.Error())
		}

		expectedLat = fmt.Sprintf("%0.6f", test.even_rLat)
		expectedLon = fmt.Sprintf("%0.6f", test.even_rLon)
		actualLat = fmt.Sprintf("%0.6f", loc.Latitude)
		actualLon = fmt.Sprintf("%0.6f", loc.Longitude)
		if (expectedLat != actualLat) {
			fmt.Errorf("Even Latitude Expected %s, got %s for test %d", expectedLat, actualLat, i)
		}
		if (expectedLon != actualLon) {
			fmt.Errorf("Even Longitude Expected %s, got %s for test %d", expectedLon, actualLon, i)
		}

		NukePlanes()
		plane = GetPlane(99887)
		plane.SetCprOddLocation(test.odd_cprLat, test.odd_cprLon, time.Now())
		plane.SetCprEvenLocation(test.even_cprLat, test.even_cprLon, time.Now())
		loc, err = plane.cprLocation.decodeSurface(test.refLat, test.refLon)

		if nil != err && test.oddErrCount == 0 {
			t.Error(err.Error())
		}

		expectedLat = fmt.Sprintf("%0.6f", test.odd_rLat)
		expectedLon = fmt.Sprintf("%0.6f", test.odd_rLon)
		actualLat = fmt.Sprintf("%0.6f", loc.Latitude)
		actualLon = fmt.Sprintf("%0.6f", loc.Longitude)

		if (expectedLat != actualLat) {
			fmt.Errorf("Odd Latitude Expected %s, got %s for test %d", expectedLat, actualLat, i)
		}
		if (expectedLon != actualLon) {
			fmt.Errorf("Odd Longitude Expected %s, got %s for test %d", expectedLon, actualLon, i)
		}
	}
}