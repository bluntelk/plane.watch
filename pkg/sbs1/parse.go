package sbs1

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	sbsMsgTypeField      = 0
	sbsMsgSubCatField    = 1
	sbsIcaoField         = 4
	sbsRecvDate          = 6
	sbsRecvTime          = 7
	sbsCallsignField     = 10
	sbsAltitudeField     = 11
	sbsGroundSpeedField  = 12
	sbsTrackField        = 13
	sbsLatField          = 14
	sbsLonField          = 15
	sbsVerticalRateField = 16
	sbsSquawkField       = 17
	sbsAlertSquawkField  = 18
	sbsEmergencyField    = 19
	sbsSpiIdentField     = 20
	sbsOnGroundField     = 21
)

type Frame struct {
	Icao         string
	Received     time.Time
	CallSign     string
	Altitude     int
	GroundSpeed  int
	Track        float64
	Lat, Lon     float64
	VerticalRate int
	Squawk       string
	Alert        string
	Emergency    string
	OnGround     bool

	HasPosition bool
}

func Parse(sbsString string) (Frame, error) {
	// decode the string
	var plane Frame
	var err error

	bits := strings.Split(sbsString, ",")
	if len(bits) != 22 {
		return plane, fmt.Errorf("Failed to Parse Input - not enough parameters: %s", sbsString)
	}

	plane.Icao = bits[sbsIcaoField]
	sTime := bits[sbsRecvDate] + " " + bits[sbsRecvTime]
	//2016/06/03 00:00:38.350
	plane.Received, err = time.Parse("2006/01/02 15:04:05.999999999", sTime)
	if nil != err {
		plane.Received = time.Now()
	}

	switch bits[sbsMsgTypeField] { // message type
	case "SEL": // SELECTION_CHANGE
		plane.CallSign = bits[sbsCallsignField]
	case "ID": // NEW_ID
		plane.CallSign = bits[sbsCallsignField]
	case "AIR": // NEW_AIRCRAFT - just indicates when a new aircraft pops up
	case "STA": // STATUS_AIRCRAFT
	// call sign field (10) contains one of:
	//	PL (Position Lost)
	// 	SL (Signal Lost)
	// 	RM (Remove)
	// 	AD (Delete)
	// 	OK (used to reset time-outs if aircraft returns into cover).
	case "CLK": // CLICK
	case "MSG": // TRANSMISSION
		switch bits[sbsMsgSubCatField] {
		case "1": // ES Identification and Category
			plane.CallSign = bits[sbsCallsignField]

		case "2": // ES Surface Position Message
			plane.Altitude, _ = strconv.Atoi(bits[sbsAltitudeField])
			plane.GroundSpeed, _ = strconv.Atoi(bits[sbsGroundSpeedField])
			plane.Track, _ = strconv.ParseFloat(bits[sbsTrackField], 32)
			plane.Lat, _ = strconv.ParseFloat(bits[sbsLatField], 32)
			plane.Lon, _ = strconv.ParseFloat(bits[sbsLonField], 32)
			plane.HasPosition = true
			plane.OnGround = "-1" == bits[sbsOnGroundField]

		case "3": // ES Airborne Position Message
			plane.Altitude, _ = strconv.Atoi(bits[sbsAltitudeField])
			plane.Lat, _ = strconv.ParseFloat(bits[sbsLatField], 32)
			plane.Lon, _ = strconv.ParseFloat(bits[sbsLonField], 32)
			plane.HasPosition = true
			plane.Alert = bits[sbsAlertSquawkField]
			plane.Emergency = bits[sbsEmergencyField]
			plane.OnGround = "-1" == bits[sbsOnGroundField]
		//SPI Flag Ignored

		case "4": // ES Airborne Velocity Message
			plane.GroundSpeed, _ = strconv.Atoi(bits[sbsGroundSpeedField])
			plane.Track, _ = strconv.ParseFloat(bits[sbsTrackField], 32)
			plane.VerticalRate, _ = strconv.Atoi(bits[sbsVerticalRateField])
			plane.OnGround = "-1" == bits[sbsOnGroundField]

		case "5": // Surveillance Alt Message
			plane.Altitude, _ = strconv.Atoi(bits[sbsAltitudeField])
			plane.Alert = bits[sbsAlertSquawkField]
			plane.OnGround = "-1" == bits[sbsOnGroundField]
			plane.CallSign = bits[sbsCallsignField]
		//SPI Flag Ignored

		case "6": // Surveillance ID Message
			plane.CallSign = bits[sbsCallsignField]
			plane.Altitude, _ = strconv.Atoi(bits[sbsAltitudeField])
			plane.Squawk = bits[sbsSquawkField]
			plane.Alert = bits[sbsAlertSquawkField]
			plane.Emergency = bits[sbsEmergencyField]
			plane.OnGround = "-1" == bits[sbsOnGroundField]
		//SPI Flag Ignored

		case "7": //Air To Air Message
			plane.Altitude, _ = strconv.Atoi(bits[sbsAltitudeField])
			plane.OnGround = "-1" == bits[sbsOnGroundField]

		case "8": // All Call Reply
			plane.OnGround = "-1" == bits[sbsOnGroundField]
		}
	}

	return plane, nil
}
