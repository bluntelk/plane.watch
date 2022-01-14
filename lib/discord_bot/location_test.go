package discord_bot

import (
	"fmt"
	"testing"
)

func TestLocationFuncs(t *testing.T) {
	var err error
	err = addAlertLocation("1", "test#1234", "test-1", 1, 1)
	if nil != err {
		t.Errorf("Should have added alert location. %s", err)
	}

	err = addAlertLocation("1", "test#1234", "test-1", 1, 1)
	if nil == err {
		t.Errorf("Should not have added same alert location twice")
	}

	err = setLocationAddress("2", "test-2", "no where")
	if nil == err {
		t.Errorf("Should not have set an address on non existant alert")
	}

	err = setLocationAddress("1", "test-1", "some place")
	if nil != err {
		t.Errorf("Should have set address on location: %s", err)
	}
	idx := getExisting("1", "test-1")
	if -1 == idx {
		t.Errorf("Should have gotten our location")
	} else {
		if "some place" != alertLocations[idx].Address {
			t.Errorf("Invalid address set")
		}
	}

	err = removeAlertLocation("1", "test-1")
	if nil != err {
		t.Errorf("Failed to remove location. %s", err)
	}

	if 0 != len(alertLocations) {
		t.Errorf("alert locations is incorrect length. expected 0, got %d", len(alertLocations))
	}
}

func TestManyAddRemove(t *testing.T) {
	var counter int
	adder := func(user, addr string) {
		counter++
		locName := fmt.Sprintf("test-loc-%d", counter)
		if err := addAlertLocation(user, user, locName, 2, 2); nil != err {
			t.Errorf("Expected to add an alert location: %s", err)
		}
		if err := setLocationAddress(user, locName, addr); nil != err {
			t.Errorf("Failed to set location address: %s", err)
		}

		if len(alertLocations) != counter {
			t.Errorf("Incorrect number of elements in alert locations array")
		}
	}

	remover := func(user string, locNum int) {
		locName := fmt.Sprintf("test-loc-%d", locNum)
		if err := removeAlertLocation(user, locName); nil != err {
			t.Errorf("Failed to remove '%s' for user '%s': %s", locName, user, err)
		}

		counter--
		if len(alertLocations) != counter {
			t.Errorf("Incorrect number of elements in alert locations array")
		}
	}

	adder("user-1", "User-1 Place 1")
	adder("user-1", "User-1 Place 2")
	adder("user-1", "User-1 Place 3")
	adder("user-2", "User-2 Place 1")
	adder("user-2", "User-2 Place 2")
	adder("user-3", "User-3 Place 1")

	remover("user-1", 1)
	remover("user-1", 2)
	remover("user-1", 3)
	remover("user-2", 4)
	remover("user-2", 5)
	remover("user-3", 6)

}
