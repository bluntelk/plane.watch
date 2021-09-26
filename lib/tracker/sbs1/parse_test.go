package sbs1

import "testing"

func TestIcaoStringToInt(t *testing.T) {
	sut := "7C1BE8"
	expected := uint32(8133608)
	icaoAddr, err := icaoStringToInt(sut)
	if err != nil {
		t.Error(err)
	}
	if icaoAddr != expected {
		t.Errorf("Expected %s to decode to %d, but got %d", sut, expected, icaoAddr)
	}
}
