package tracker

import (
	"testing"
	"time"
	"fmt"
)

func TestNLFunc(t *testing.T) {
	for i, f := range NLTable {
		if r := getNumLongitudeZone(f - 0.01); i != r {
			t.Errorf("NL Table Fail: Expected %0.2f to yield %d, got %d", f, i, r)
		}
	}
}

func TestCprDecode(t *testing.T) {
	type testDataType struct {
		evenLat, evenLon float64
		oddLat, oddLon   float64


		evenRlatCheck1, evenRlonCheck1 string

		evenRlat, evenRlon string
		oddRlat, oddRlon   string
	}
	testData := []testDataType{
		//odd *8d7c4516581f76e48d95e8ab20ca; even *8d7c4516581f6288f83ade534ae1;
		{evenLat: 83068, evenLon:15070, oddLat:94790, oddLon:103912, oddRlat:"-32.197483", oddRlon:"+116.028629", evenRlat:"-32.197449", evenRlon:"+116.027820"},

		// odd *8d7c4516580f06fc6d8f25d8669d; even *8d7c4516580df2a168340b32212a;
		{evenLat: 86196, evenLon:13323, oddLat:97846, oddLon:102181, oddRlat:"-32.055219", oddRlon:"+115.931602", evenRlat:"-32.054260", evenRlon:"+115.931854"},

		// test data from cprtest.c from mutability dump1090
		{evenLat: 80536, evenLon:9432, oddLat:61720, oddLon:9192, evenRlat:"+51.686646", evenRlon:"+0.700156", oddRlat:"+51.686763", oddRlon:"+0.701294"},
	}
	airDlat0 := "+6.000000"
	airDlat1 := "+6.101695"
	trk := NewTracker()

	for i, d := range testData {
		plane := trk.GetPlane(11234)

		plane.SetCprOddLocation(d.oddLat, d.oddLon, time.Now())
		time.Sleep(2)
		plane.SetCprEvenLocation(d.evenLat, d.evenLon, time.Now())
		loc, err := plane.cprLocation.decodeGlobalAir()
		if err != nil {
			t.Error(err)
		}

		lat := fmt.Sprintf("%+0.6f", loc.Latitude);
		lon := fmt.Sprintf("%+0.6f", loc.Longitude);

		if lat != d.oddRlat {
			t.Errorf("Plane Latitude is wrong for packet %d: should be %s, was %s", i, d.oddRlat, lat)
		}
		if lon != d.oddRlon {
			t.Errorf("Plane Latitude is wrong for packet %d: should be %s, was %s", i, d.oddRlon, lon)
		}

		if airDlat0 != fmt.Sprintf("%+0.6f", plane.cprLocation.airDLat0) {
			t.Error("AirDlat0 is wrong")
		}
		if airDlat1 != fmt.Sprintf("%+0.6f", plane.cprLocation.airDLat1) {
			t.Error("AirDlat1 is wrong")
		}

		plane.SetCprEvenLocation(d.evenLat, d.evenLon, time.Now())
		time.Sleep(2)
		plane.SetCprOddLocation(d.oddLat, d.oddLon, time.Now())
		loc, err = plane.cprLocation.decodeGlobalAir()
		if err != nil {
			t.Error(err)
		}

		lat = fmt.Sprintf("%+0.6f", loc.Latitude);
		lon = fmt.Sprintf("%+0.6f", loc.Longitude);

		if lat != d.evenRlat {
			t.Errorf("Plane Latitude is wrong for packet %d: should be %s, was %s", i, d.evenRlat, lat)
		}
		if lon != d.evenRlon {
			t.Errorf("Plane Latitude is wrong for packet %d: should be %s, was %s", i, d.evenRlon, lon)
		}

		if airDlat0 != fmt.Sprintf("%+0.6f", plane.cprLocation.airDLat0) {
			t.Error("AirDlat0 is wrong")
		}
		if airDlat1 != fmt.Sprintf("%+0.6f", plane.cprLocation.airDLat1) {
			t.Error("AirDlat1 is wrong")
		}

	}
}
