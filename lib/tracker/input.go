package tracker

import (
	"errors"
	"plane.watch/lib/tracker/mode_s"
	sbs12 "plane.watch/lib/tracker/sbs1"
	"time"
)

func HandleAvr(avr string, received time.Time) error {

	f, err := mode_s.DecodeString(avr, received)
	if nil != err {
		return err
	}
	HandleModeSFrame(f)

	return nil
}

func HandleSbs1(sbs1 string, received time.Time) error {
	f, err := sbs12.Parse(sbs1)
	if nil != err {
		return err
	}
	if nil != HandleSbs1Frame(f) {
		return errors.New("failed to understand SBS1 information")
	}
	return nil
}