package calculation

import (
	"context"
	"log"
	"strconv"
	"time"

	"googlemaps.github.io/maps"
)

var mapsClient *maps.Client
var startTimeForRoutes time.Time
var numberOfRoutes int64

// FindBestRoute finds the fastest route to traverse every stop, every stop must have an edge to every other stop
func FindBestRoute(stops []string, startTime time.Time) ([]string, time.Duration) {
	if mapsClient == nil {
		var err error
		mapsClient, err = maps.NewClient(maps.WithAPIKey("Totally Real API Key"))
		if err != nil {
			log.Fatalf("fatal error: %s", err)
		}
	}

	numberOfRoutes = 0
	startTimeForRoutes = startTime
	route := make([]string, 0)
	return findBestRouteHelper(route, stops)
}

func findBestRouteHelper(curRoute []string, stopsLeft []string) ([]string, time.Duration) {
	if len(stopsLeft) == 0 {
		numberOfRoutes++
		return curRoute, findRouteTime(curRoute)
	}

	var bestRoute []string
	var bestDuration time.Duration

	for i := range stopsLeft {
		route, duration := findBestRouteHelper(append(curRoute, stopsLeft[i]), removeIndex(i, stopsLeft))
		if duration < bestDuration {
			bestRoute = route
			bestDuration = duration
		}
	}

	return bestRoute, bestDuration
}

func removeIndex(index int, list []string) []string {
	list[index] = list[len(list)-1]
	return list[:len(list)-1]
}

func findRouteTime(route []string) time.Duration {
	startTime := startTimeForRoutes.UTC().Second()
	var duration time.Duration
	for i := 0; i < len(route)-1; i++ {
		duration += findEdgeTime(route[i], route[i+1], startTime+int(duration))
	}
	return duration
}

func findEdgeTime(stopA string, stopB string, startTime int) time.Duration {
	req := &maps.DistanceMatrixRequest{
		Origins:       []string{stopA},
		Destinations:  []string{stopB},
		DepartureTime: strconv.Itoa(startTime),
		Mode:          maps.TravelModeTransit,
		TransitMode: []maps.TransitMode{
			maps.TransitModeRail,
			maps.TransitModeSubway,
			maps.TransitModeTrain,
			maps.TransitModeTram,
		},
	}

	resp, err := mapsClient.DistanceMatrix(context.Background(), req)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	duration := resp.Rows[0].Elements[0].Duration
	return duration
}
