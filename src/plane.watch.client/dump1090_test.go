package main

import (
	"testing"
	"mode_s"
	"time"
	"tracker"
	"fmt"
	"os"
)

func TestTracking(t *testing.T) {
	frames := []string{
		"*8D40621D58C382D690C8AC2863A7;",
		"*8D40621D58C386435CC412692AD6;",
	}
	performTrackingTest(frames, t)

	plane := tracker.GetPlane(4219421)
	if plane.Location.Altitude != 38000 {
		t.Error("Plane should be at 38000 feet")
	}

	lat := "+52.2572021484375"
	lon := "+3.9193725585938";
	if lon != fmt.Sprintf("%+03.13f", plane.Location.Longitude) {
		t.Errorf("Longitude Calculation was incorrect: expected %s, got %+0.13f", lon, plane.Location.Longitude)
	}
	if lat != fmt.Sprintf("%+03.13f", plane.Location.Latitude) {
		t.Errorf("Latitude Calculation was incorrect: expected %s, got %+0.13f", lat, plane.Location.Latitude)
	}
}

func TestTracking2(t *testing.T) {
	frames := []string{
		"*8D7C7DAA99146D0980080D6131A1;",
		"*5D7C7DAACD3CE9;",
		"*0005050870B303;",
		"*8D7C7DAA99146C0980040D2A616F;",
		"*8D7C7DAAF80020060049B06CA244;",
		"*8D7C7DAA582886FA618B21ADB377;",
		"*5D7C7DAACD3CE9;",
		"*8D7C7DAA5828829F322FE81F6DD1;",
		"*8D7C7DAA99146C0980040D2A616F;",
		"*8D7C7DAA99146C0980040D2A616F;",
		"*8D7C7DAA99146C0960080D47BBB9;",
		"*8D7C7DAA582886FA778B115D2F89;",
		"*000005084A3646;",
		"*000005084A3646;",
		"*28000A00307264;",
		"*8D7C7DAA99146A09280C0D91E947;",
		"*8D7C7DAA9914690920080DC2621D;",
		"*8D7C7DAA9914690928040DE49A15;",
		"*8D7C7DAA210DA1E0820820472D63;",
		"*5D7C7DAACD3CE9;",
		"*8D7C7DAA582886FB218A9AFB0420;",
		"*5D7C7DAACD3CE9;",
		"*8D7C7DAA5828829FF42F5E556B2D;",
		"*8D7C7DAA9914680920080DC168D3;",
		"*000005084A3646;",
		"*5D7C7DAACD3CE9;",
		"*8D7C7DAA582886FB318A8FD96CD7;",
		"*8D7C7DAA9914670900080D9576E0;",
		"*000005084A3646;",
	}
	performTrackingTest(frames, t)
}

func performTrackingTest(frames []string, t *testing.T) {
	f, _ := os.Open(os.DevNull)
	tracker.SetDebugOutput(f)
	for _, msg := range frames {
		frame, err := mode_s.DecodeString(msg, time.Now())
		if nil != err {
			t.Errorf("%s", err)
		}
		tracker.HandleModeSFrame(frame, false)
	}
}

func TestAltitudeDecode(t *testing.T) {
	frame, err := mode_s.DecodeString("*8D7C7DAA582886FB218A9AFB0420;", time.Now())
	if nil != err {
		t.Error(err)
	}
	if 600 != frame.Altitude() {
		t.Errorf("Expected an altitude of 600 feet, got %d", frame.Altitude())
	}
}

func TestCprDecode(t *testing.T) {
	type testDataType struct {
		raw      string
		icoa     string
		isEven   bool
		alt      int32
		lat, lon string
	}
	testData := []testDataType{
		{raw: "*8d7c4516581f76e48d95e8ab20ca;", icoa:"7c4516", isEven: false, alt:5175, lat:"+0.000000", lon:"+0.000000"},
		{raw: "*8d7c4516581f6288f83ade534ae1;", icoa:"7c4516", isEven: true, alt:5150, lat:"-32.197483", lon:"+116.028629"},

		//{raw: "*8d7c4516580f06fc6d8f25d8669d;", icoa:"7c4516", isEven: false, alt:1800, lat:"-0.000000", lon:"+0.000000"},
		//{raw: "*8d7c4516580df2a168340b32212a;", icoa:"7c4516", isEven: true, alt:1775, lat:"-32.054260", lon:"+115.931854"},
	}

	for i, d := range testData {
		time.Sleep(1)
		frame, err := mode_s.DecodeString(d.raw, time.Now())
		if nil != err {
			t.Error(err)
		}
		plane := tracker.HandleModeSFrame(frame, false)

		if nil == plane {
			t.Errorf("Plane data should have been updated")
		}

		if plane.Location.Altitude != d.alt {
			t.Errorf("Plane Altitude is wrong for packet %d: should be %d, was %d", i, d.alt, plane.Location.Altitude)
		}

		lat := fmt.Sprintf("%+0.6f", plane.Location.Latitude);
		lon := fmt.Sprintf("%+0.6f", plane.Location.Longitude);

		if lat != d.lat {
			t.Errorf("Plane Latitude is wrong for packet %d: should be %s was %s", i, d.lat, lat)
		}
		if lon != d.lon {
			t.Errorf("Plane Latitude is wrong for packet %d: should be %s was %s", i, d.lon, lon)
		}
	}
}