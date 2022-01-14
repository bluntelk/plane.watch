package mapping

import (
	"context"
	"errors"
	"fmt"
	"log"
	"plane.watch/lib/mapping/here"
	"strings"
)

var hereMapsApi = ""

func SetHereMapsApiKey(apiKey string) {
	hereMapsApi = apiKey
}

// FindCoordinates uses whatever tool we have configured to do the lookup
func FindCoordinates(location string) (float64, float64, error) {
	if "" != hereMapsApi {
		return findCoordinatesHereMaps(location)
	}
	return 0, 0, errors.New("you do not have Here Maps API provided")
}

func findCoordinatesHereMaps(location string) (float64, float64, error) {
	cfg := here.NewConfiguration()

	client := here.NewAPIClient(cfg)

	apiKey := map[string]here.APIKey{
		"ApiKey": {
			Key:    hereMapsApi,
			Prefix: "",
		},
	}
	ctx := context.WithValue(context.Background(), here.ContextAPIKeys, apiKey)
	resp, httpResp, err := client.DefaultApi.
		GeocodeGet(ctx).
		Q(location).
		Limit(5).
		Lang([]string{"en"}).
		Execute()

	if "" != err.Error() || httpResp.StatusCode != 200 {
		log.Printf("Failed looking up %s: [HTTP:%d] %s", location, httpResp.StatusCode, err)
		return 0, 0, errors.New("failed to lookup this location")
	}

	if 0 == len(resp.Items) {
		log.Printf("Failed looking up location %s, no results", location)
		return 0, 0, errors.New("did not find that location")
	}

	if 1 == len(resp.Items) {
		lat, okLat := resp.Items[0].Position.GetLatOk()
		lon, okLon := resp.Items[0].Position.GetLngOk()
		if !okLat || !okLon {
			log.Printf("Returned result did not have a valid lat/lon. %+v", resp.Items[0])
			return 0, 0, errors.New("found location did not have valid coords")
		}
		return *lat, *lon, nil
	}

	log.Printf("Found many addresses for lookup")
	// we have >1 results, let's combine them into a single error for the user to choose from
	var buf strings.Builder
	buf.WriteString("Multiple results were returned, please use a more specific value.")

	for _, item := range resp.Items {
		buf.WriteString(fmt.Sprintf("```%s``` coords: %0.5f,%0.5f", item.GetTitle(), item.Position.GetLat(), item.Position.GetLng()))
	}

	return 0, 0, errors.New(buf.String())
}
