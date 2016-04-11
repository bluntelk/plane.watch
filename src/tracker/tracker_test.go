package tracker

import "testing"

func TestNLFunc(t *testing.T) {
	for i, f := range NLTable {
		if r := getNumLongitudeZone(f-0.01); i != r {
			t.Errorf("NL Table Fail: Expected %0.2f to yield %d, got %d", f, i, r)
		}
	}
}