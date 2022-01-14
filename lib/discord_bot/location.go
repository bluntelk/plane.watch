package discord_bot

// handles the list of alerting locations

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"plane.watch/lib/mapping"
	"strings"
	"sync"
)

const alertLocationsFile = "alert-locations.json"

type (
	location struct {
		DiscordUserId   string
		DiscordUserName string
		LocationName    string
		Address         string
		Lat             float64
		Lon             float64
		AlertRadius     int // The radius of the circle to alert in
	}
)

var (
	alertLocationsRWLock sync.RWMutex
	alertLocations       []location
	isLoaded             bool
)

func getPath() string {
	binaryPath, _ := os.Executable()
	if "" != binaryPath && !strings.Contains(binaryPath, "/go-build/") {
		return path.Dir(binaryPath)
	}
	dir, _ := os.Getwd()
	return dir
}

// getExisting returns the id in the array of the existing record, -1 if it does not exist
func getExisting(discordUserId, locationName string) int {
	alertLocationsRWLock.RLock()
	defer alertLocationsRWLock.RUnlock()
	for idx, loc := range alertLocations {
		if loc.DiscordUserId != discordUserId {
			continue
		}
		if loc.LocationName == locationName {
			return idx
		}
	}
	return -1
}

func getLocationsForUser(discordUserId string) []location {
	alertLocationsRWLock.RLock()
	defer alertLocationsRWLock.RUnlock()

	var ret []location
	for _, loc := range alertLocations {
		if loc.DiscordUserId == discordUserId {
			ret = append(ret, loc)
		}
	}
	return ret
}

func addAlertLocation(discordUserId, discordUserName, locationName string, lat, lon float64) error {
	// make sure we do not already have this name
	if -1 != getExisting(discordUserId, locationName) {
		return errors.New("you have an existing location with this name")
	}

	alertLocationsRWLock.Lock()
	loc := location{
		DiscordUserId:   discordUserId,
		DiscordUserName: discordUserName,
		LocationName:    locationName,
		Lat:             lat,
		Lon:             lon,
		AlertRadius:     500,
	}
	alertLocations = append(alertLocations, loc)
	alertLocationsRWLock.Unlock()
	return saveLocationsList()
}

func removeAlertLocation(discordUserId, locationName string) error {
	idx := getExisting(discordUserId, locationName)
	if -1 == idx {
		return errors.New("unknown location")
	}
	alertLocationsRWLock.Lock()
	if len(alertLocations) == 1 && idx == 0 {
		alertLocations = []location{}
	} else if idx == len(alertLocations)-1 {
		// last element, just shorten
		alertLocations = alertLocations[:idx-1]
	} else {
		alertLocations = append(alertLocations[:idx], alertLocations[idx+1:]...)
	}
	alertLocationsRWLock.Unlock()
	return saveLocationsList()
}

// setLocationAddress allows us to set the address of a geocoded location
func setLocationAddress(discordUserId, locationName, address string) error {
	idx := getExisting(discordUserId, locationName)
	if -1 == idx {
		return errors.New("that location name does not exist")
	}
	alertLocationsRWLock.Lock()
	alertLocations[idx].Address = address
	alertLocationsRWLock.Unlock()
	return saveLocationsList()
}

// allows updating the radius in which we alert for this location
func setLocationAlertRadius(discordUserId, locationName string, alertRadius int) error {
	idx := getExisting(discordUserId, locationName)
	if -1 == idx {
		return errors.New("that location name does not exist")
	}
	alertLocationsRWLock.Lock()
	alertLocations[idx].AlertRadius = alertRadius
	alertLocationsRWLock.Unlock()
	return saveLocationsList()
}

func loadLocationsList() {
	alertLocationsRWLock.Lock()
	defer alertLocationsRWLock.Unlock()
	if isLoaded {
		return
	}
	saveLoc := getPath() + "/" + alertLocationsFile
	b, err := ioutil.ReadFile(saveLoc)
	if nil != err {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("No save file. %s does not exist. proceeding with empty list", saveLoc)
			return
		}
		log.Fatalf("Failed to read %s. %s", saveLoc, err)
		return
	}
	err = json.Unmarshal(b, &alertLocations)
	if nil != err {
		log.Fatalf("Failed to parse %s JSON perfectly. %s", saveLoc, err)
	}
	isLoaded = true
}

func saveLocationsList() error {
	alertLocationsRWLock.RLock()
	defer alertLocationsRWLock.RUnlock()
	saveLoc := getPath() + "/" + alertLocationsFile

	b, err := json.MarshalIndent(alertLocations, "", "  ")
	if nil != err {
		if len(b) <= 3 {
			return fmt.Errorf("we don't have a good marshalling, not saving. %s", err)
		} else {
			log.Printf("Failed to marshal JSON, attempting to save what we have. %s", err)
		}
	}

	err = ioutil.WriteFile(saveLoc, b, 0644)
	if nil != err {
		return fmt.Errorf("failed to save locations to %s. %s", saveLoc, err)
	}
	return nil
}

func geoCodeAddress(addr string) (float64, float64, error) {
	log.Printf("Geocoding user address: %s", addr)
	return mapping.FindCoordinates(addr)
}
