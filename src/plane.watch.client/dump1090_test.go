package main

import (
	"testing"
	"tracker"
	"mode_s"
	"time"
	"fmt"
)

func TestTracking(t *testing.T) {
	frames := []string{
		"*8D40621D58C382D690C8AC2863A7;",
		"*8D40621D58C386435CC412692AD6;",
	}

	for _, msg := range frames {
		frame, err := mode_s.DecodeString(msg, time.Now())
		if nil != err {
			t.Errorf("%s", err)
		}
		tracker.HandleModeSFrame(frame)
	}

	plane := tracker.GetPlane(4219421)
	if plane.Altitude != 38000 {
		t.Error("Plane should be at 38000 feet")
	}

	lat := "+52.2657801741261"
	lon := "+3.9389125279018";
	if lon != fmt.Sprintf("%+03.13f", plane.Location.Longitude) {
		t.Errorf("Longitude Calculation was incorrect: expected %s, got %0.12f", lon, plane.Location.Longitude)
	}
	if lat != fmt.Sprintf("%+03.13f", plane.Location.Latitude) {
		t.Errorf("Latitude Calculation was incorrect: expected %s, got %0.13f", lat, plane.Location.Latitude)
	}

}
