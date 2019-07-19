package datacollection

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/mellena1/mbta-v3-go/mbta"
)

// GetStopCoordinatesForGoogleAPI calculates and saves mbta api location data
func GetStopCoordinatesForGoogleAPI() {
	mbtaClient := mbta.NewClient(mbta.ClientConfig{
		APIKey: "",
	})

	endpoints, err := ImportStopsFromFile(Endpoints)
	if err != nil {
		log.Fatalf("Failed to import endpoint list: %s", err)
	}

	reqConfig := mbta.GetAllStopsRequestConfig{}
	stops, _, err := mbtaClient.Stops.GetAllStops(&reqConfig)
	if err != nil {
		log.Fatalf("Failed to get stops from mbta API: %s", err)
	}

	endpointNameStopMap := make(map[string]*Stop, len(endpoints))
	for _, endpoint := range endpoints {
		endpointNameStopMap[endpoint.Name] = endpoint
	}

	for _, stop := range stops {
		if stop.ParentStation != nil {
			continue
		}
		if endpointStop, ok := endpointNameStopMap[stop.Name]; ok {
			endpointStop.SetLongitudeCommaLatitude(stop.Longitude, stop.Latitude)
		}
	}

	err = ExportLocationData(endpoints)
	if err != nil {
		log.Fatalf("Failed to export stop location data: %s", err)
	}
}

// ExportLocationData exorts the location data from the mbta API
func ExportLocationData(stops []*Stop) error {
	data, err := json.Marshal(stops)
	if err != nil {
		return err
	}

	ioutil.WriteFile(string(StopLocations), data, 0644)

	return nil
}
