package sbs1

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	sbs_msg_type_field      = 0
	sbs_msg_sub_cat_field   = 1
	sbs_icao_field          = 4
	sbs_recv_date           = 6
	sbs_recv_time           = 7
	sbs_callsign_field      = 10
	sbs_altitude_field      = 11
	sbs_ground_speed_field  = 12
	sbs_track_field         = 13
	sbs_lat_field           = 14
	sbs_lon_field           = 15
	sbs_vertical_rate_field = 16
	sbs_squawk_field        = 17
	sbs_alert_squawk_field  = 18
	sbs_emergency_field     = 19
	sbs_spi_ident_field     = 20
	sbs_on_ground_field     = 21
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

	plane.Icao = bits[sbs_icao_field]
	sTime := bits[sbs_recv_date] + " " + bits[sbs_recv_time]
	//2016/06/03 00:00:38.350
	plane.Received, err = time.Parse("2006/01/02 15:04:05.999999999", sTime)
	if nil != err {
		plane.Received = time.Now()
	}

	switch bits[sbs_msg_type_field] { // message type
	case "SEL": // SELECTION_CHANGE
		plane.CallSign = bits[sbs_callsign_field]
	case "ID": // NEW_ID
		plane.CallSign = bits[sbs_callsign_field]
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
		switch bits[sbs_msg_sub_cat_field] {
		case "1": // ES Identification and Category
			plane.CallSign = bits[sbs_callsign_field]

		case "2": // ES Surface Position Message
			plane.Altitude, _ = strconv.Atoi(bits[sbs_altitude_field])
			plane.GroundSpeed, _ = strconv.Atoi(bits[sbs_ground_speed_field])
			plane.Track, _ = strconv.ParseFloat(bits[sbs_track_field], 32)
			plane.Lat, _ = strconv.ParseFloat(bits[sbs_lat_field], 32)
			plane.Lon, _ = strconv.ParseFloat(bits[sbs_lon_field], 32)
			plane.HasPosition = true
			plane.OnGround = "-1" == bits[sbs_on_ground_field]

		case "3": // ES Airborne Position Message
			plane.Altitude, _ = strconv.Atoi(bits[sbs_altitude_field])
			plane.Lat, _ = strconv.ParseFloat(bits[sbs_lat_field], 32)
			plane.Lon, _ = strconv.ParseFloat(bits[sbs_lon_field], 32)
			plane.HasPosition = true
			plane.Alert = bits[sbs_alert_squawk_field]
			plane.Emergency = bits[sbs_emergency_field]
			plane.OnGround = "-1" == bits[sbs_on_ground_field]
		//SPI Flag Ignored

		case "4": // ES Airborne Velocity Message
			plane.GroundSpeed, _ = strconv.Atoi(bits[sbs_ground_speed_field])
			plane.Track, _ = strconv.ParseFloat(bits[sbs_track_field], 32)
			plane.VerticalRate, _ = strconv.Atoi(bits[sbs_vertical_rate_field])
			plane.OnGround = "-1" == bits[sbs_on_ground_field]

		case "5": // Surveillance Alt Message
			plane.Altitude, _ = strconv.Atoi(bits[sbs_altitude_field])
			plane.Alert = bits[sbs_alert_squawk_field]
			plane.OnGround = "-1" == bits[sbs_on_ground_field]
			plane.CallSign = bits[sbs_callsign_field]
		//SPI Flag Ignored

		case "6": // Surveillance ID Message
			plane.CallSign = bits[sbs_callsign_field]
			plane.Altitude, _ = strconv.Atoi(bits[sbs_altitude_field])
			plane.Squawk = bits[sbs_squawk_field]
			plane.Alert = bits[sbs_alert_squawk_field]
			plane.Emergency = bits[sbs_emergency_field]
			plane.OnGround = "-1" == bits[sbs_on_ground_field]
		//SPI Flag Ignored

		case "7": //Air To Air Message
			plane.Altitude, _ = strconv.Atoi(bits[sbs_altitude_field])
			plane.OnGround = "-1" == bits[sbs_on_ground_field]

		case "8": // All Call Reply
			plane.OnGround = "-1" == bits[sbs_on_ground_field]
		}
	}

	return plane, nil
}
