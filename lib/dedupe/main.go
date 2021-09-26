package dedupe

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/beast"
	"plane.watch/lib/tracker/mode_s"
	"plane.watch/lib/tracker/sbs1"
	"sync"
	"time"
)

/**
This package provides a way to deduplicate mode_s messages.

Consider a message a duplicate if we have seen it in the last minute
 */

type (
	ForgetfulSynMap struct {
		lookup *sync.Map
		sweeper  *time.Timer
		sweepInterval time.Duration
		oldAfter time.Duration
	}

	Filter struct {
		list *ForgetfulSynMap
	}
)

func NewFilter() *Filter {
	return &Filter{
		list: NewForgetfulSyncMap(),
	}
}

func NewForgetfulSyncMap() *ForgetfulSynMap {
	f := ForgetfulSynMap{
		lookup: &sync.Map{},
		sweepInterval: time.Second * 10,
		oldAfter: time.Minute,
	}
	f.sweeper = time.AfterFunc(f.oldAfter, func() {
		f.sweep()
		f.sweeper.Reset(f.sweepInterval)
	})

	return &f
}

func (f *ForgetfulSynMap) sweep() {
	var remove bool
	removeCount := 0
	testCount := 0
	oldest := time.Now().Add(-time.Minute)
	f.lookup.Range(func(key, value interface{}) bool {
		remove = true
		testCount++
		if t, ok := value.(time.Time); ok {
			if t.After(oldest) {
				remove = false
			}
		}

		if remove {
			f.lookup.Delete(key)
			removeCount++
		}

		return true
	})
	log.Debug().Msgf("Removed %d old of %d entries", removeCount, testCount)
}

func (f *ForgetfulSynMap) HasKey(key interface{}) bool {
	if _, ok :=f.lookup.Load(key); ok {
		return true
	}
	return false
}

func (f *ForgetfulSynMap) AddKey(key interface{}) {
	// avoid storing empty things
	if nil == key {return}
	if kb, ok := key.([]byte); ok {
		if 0 == len(kb) {
			return
		}
	}
	if ks, ok := key.(string); ok {
		if "" == ks {
			return
		}
	}
	f.lookup.Store(key, time.Now())
}

func (f *Filter) DeDupe(frame tracker.Frame) tracker.Frame {
	if nil == frame {
		return nil
	}
	var key interface{}
	switch (frame).(type) {
	case *beast.Frame:
		key = fmt.Sprintf("%X", frame.(*beast.Frame).AvrRaw())
	case *mode_s.Frame:
		key = fmt.Sprintf("%X", frame.(*mode_s.Frame).Raw())
	case *sbs1.Frame:
		// todo: investigate better dedupe detection for sbs1
		key = string(frame.(*sbs1.Frame).Raw())
	default:
	}
	if f.list.HasKey(key) {
		return nil
	}
	f.list.AddKey(key)
	return frame
}
