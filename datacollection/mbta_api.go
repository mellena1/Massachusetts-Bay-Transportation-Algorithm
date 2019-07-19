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

	reqConfig := mbta.GetAllStopsRequestConfig{}
	stops, _, err := mbtaClient.Stops.GetAllStops(&reqConfig)
	if err != nil {
		log.Fatalf("Failed to get stops from mbta API: %s", err)
	}

	specialEdges, err := ReadSpecialEdgesFromFile(SpecialEdgesFile)
	if err != nil {
		log.Fatalf("Failed to import special edges data: %s", err)
	}

	endpoints, err := ImportStopsFromFile(Endpoints)
	if err != nil {
		log.Fatalf("Failed to import endpoint list: %s", err)
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
			continue
		}
		for _, middleStop := range specialEdges {
			if stop.Name == middleStop.Name {
				middleStop.SetLongitudeCommaLatitude(stop.Longitude, stop.Latitude)
			}
		}
	}

	err = ExportSpecialEdgesLocationData(specialEdges)
	if err != nil {
		log.Fatalf("Failed to export stop location data: %s", err)
	}

	err = ExportLocationData(endpoints)
	if err != nil {
		log.Fatalf("Failed to export stop location data: %s", err)
	}
}

// ExportSpecialEdgesLocationData exorts the location data from the mbta API for the middle stops of special edges
func ExportSpecialEdgesLocationData(specialEdges SpecialEdges) error {
	data, err := json.Marshal(specialEdges)
	if err != nil {
		return err
	}

	ioutil.WriteFile(string(SpecialEdgesFileWithLocationData), data, 0644)

	return nil
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
