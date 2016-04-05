package mode_s

import (
	"testing"
)

func TestAC12Decode(t *testing.T) {
	doAc12Decode(255, 2175, t) // Q Bit Set
	doAc12Decode(648, 999, t)  // Q Bit Clear
}

func doAc12Decode(alt, expected int32, t *testing.T) {
	tested := decodeAC12Field(alt)
	if tested != expected {
		t.Errorf("Expected %d feet for input %d but got %d instead", expected, alt, tested)
	}
}