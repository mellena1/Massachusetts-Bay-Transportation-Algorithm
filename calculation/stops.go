package calculation

import (
	"log"
	"strconv"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/graph"
	"github.com/mellena1/mbta-v3-go/mbta"
)

// Stop stores the coordinates of a stop for google maps
type Stop struct {
	Name      string
	Longitude string
	Latitude  string
}

func (s *Stop) getCoordinateString() string {
	return s.Latitude + "," + s.Longitude
}

var mbtaClient *mbta.Client

// GetEndpointStops returns a list of endpoint stops
func GetEndpointStops() []Stop {
	if mbtaClient == nil {
		mbtaClient = mbta.NewClient(mbta.ClientConfig{
			APIKey: "",
		})
	}

	reqConfig := mbta.GetAllStopsRequestConfig{}
	stops, _, err := mbtaClient.Stops.GetAllStops(&reqConfig)
	if err != nil {
		log.Fatalf("A fatal error occurred: %s", err)
	}

	_, jsonStopIDs, err := graph.LoadGraphFile("graph.json")
	if err != nil {
		log.Fatalf("A fatal error occurred: %s", err)
	}

	endpoints := make([]Stop, 0)
	for _, stop := range stops {
		// only care about parent stations
		if stop.ParentStation != nil {
			continue
		}
		if v, ok := jsonStopIDs[stop.ID]; ok {
			if v.IsEndpoint() {
				endpoints = append(endpoints, Stop{
					Name:      stop.Name,
					Longitude: strconv.FormatFloat(stop.Longitude, 'f', -1, 64),
					Latitude:  strconv.FormatFloat(stop.Latitude, 'f', -1, 64),
				})
			}
		}
	}

	return endpoints
}
