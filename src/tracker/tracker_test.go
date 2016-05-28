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
		even_lat, even_lon                 float64
		odd_lat, odd_lon                   float64


		even_rlat_check1, even_rlon_check1 string

		even_rlat, even_rlon               string
		odd_rlat, odd_rlon                 string
	}
	testData := []testDataType{
		//odd *8d7c4516581f76e48d95e8ab20ca; even *8d7c4516581f6288f83ade534ae1;
		{even_lat:83068, even_lon:15070, odd_lat:94790, odd_lon:103912, odd_rlat:"-32.197483", odd_rlon:"+116.028629", even_rlat:"-32.197449", even_rlon:"+116.027820"},

		// odd *8d7c4516580f06fc6d8f25d8669d; even *8d7c4516580df2a168340b32212a;
		{even_lat:86196, even_lon:13323, odd_lat:97846, odd_lon:102181, odd_rlat:"-32.055219", odd_rlon:"+115.931602", even_rlat:"-32.054260", even_rlon:"+115.931854"},

		// test data from cprtest.c from mutability dump1090
		{even_lat:80536, even_lon:9432, odd_lat:61720, odd_lon:9192, even_rlat:"+51.686646", even_rlon:"+0.700156", odd_rlat:"+51.686763", odd_rlon:"+0.701294"},
	}
	airDlat0 := "+6.000000";
	airDlat1 := "+6.101695";

	for i, d := range testData {
		plane := GetPlane(11234)

		plane.SetCprOddLocation(d.odd_lat, d.odd_lon, time.Now())
		time.Sleep(2)
		plane.SetCprEvenLocation(d.even_lat, d.even_lon, time.Now())
		loc, err := plane.cprLocation.decodeAir()
		if err != nil {
			t.Error(err)
		}

		lat := fmt.Sprintf("%+0.6f", loc.Latitude);
		lon := fmt.Sprintf("%+0.6f", loc.Longitude);

		if lat != d.odd_rlat {
			t.Errorf("Plane Latitude is wrong for packet %d: should be %s, was %s", i, d.odd_rlat, lat)
		}
		if lon != d.odd_rlon {
			t.Errorf("Plane Latitude is wrong for packet %d: should be %s, was %s", i, d.odd_rlon, lon)
		}

		if airDlat0 != fmt.Sprintf("%+0.6f", plane.cprLocation.airDLat0) {
			t.Error("AirDlat0 is wrong")
		}
		if airDlat1 != fmt.Sprintf("%+0.6f", plane.cprLocation.airDLat1) {
			t.Error("AirDlat1 is wrong")
		}

		plane.SetCprEvenLocation(d.even_lat, d.even_lon, time.Now())
		time.Sleep(2)
		plane.SetCprOddLocation(d.odd_lat, d.odd_lon, time.Now())
		loc, err = plane.cprLocation.decodeAir()
		if err != nil {
			t.Error(err)
		}

		lat = fmt.Sprintf("%+0.6f", loc.Latitude);
		lon = fmt.Sprintf("%+0.6f", loc.Longitude);

		if lat != d.even_rlat {
			t.Errorf("Plane Latitude is wrong for packet %d: should be %s, was %s", i, d.even_rlat, lat)
		}
		if lon != d.even_rlon {
			t.Errorf("Plane Latitude is wrong for packet %d: should be %s, was %s", i, d.even_rlon, lon)
		}

		if airDlat0 != fmt.Sprintf("%+0.6f", plane.cprLocation.airDLat0) {
			t.Error("AirDlat0 is wrong")
		}
		if airDlat1 != fmt.Sprintf("%+0.6f", plane.cprLocation.airDLat1) {
			t.Error("AirDlat1 is wrong")
		}

	}
}