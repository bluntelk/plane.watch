package dedupe

import (
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type (
	ForgetfulSyncMap struct {
		lookup        *sync.Map
		sweeper       *time.Timer
		sweepInterval time.Duration
		oldAfter      time.Duration
	}
)

func NewForgetfulSyncMap() *ForgetfulSyncMap {
	f := ForgetfulSyncMap{
		lookup:        &sync.Map{},
		sweepInterval: time.Second * 10,
		oldAfter:      time.Minute,
	}
	f.sweeper = time.AfterFunc(f.oldAfter, func() {
		f.sweep()
		f.sweeper.Reset(f.sweepInterval)
	})

	return &f
}

func (f *ForgetfulSyncMap) sweep() {
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

func (f *ForgetfulSyncMap) HasKey(key interface{}) bool {
	if _, ok := f.lookup.Load(key); ok {
		return true
	}
	return false
}

func (f *ForgetfulSyncMap) AddKey(key interface{}) {
	// avoid storing empty things
	if nil == key {
		return
	}
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
