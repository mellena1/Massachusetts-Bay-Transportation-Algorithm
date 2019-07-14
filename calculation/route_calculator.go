package calculation

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"googlemaps.github.io/maps"
)

var mapsClient *maps.Client
var startTimeForRoutes time.Time
var numberOfRoutes int64

// FindBestRoute finds the fastest route to traverse every stop, every stop must have an edge to every other stop
func FindBestRoute(stops []Stop, startTime time.Time) ([]Stop, time.Duration) {
	if mapsClient == nil {
		var err error
		mapsClient, err = maps.NewClient(maps.WithAPIKey("put api key here"))
		if err != nil {
			log.Fatalf("fatal error: %s", err)
		}
	}

	numberOfRoutes = 0
	startTimeForRoutes = startTime
	route := make([]Stop, 0)
	return findBestRouteHelper(route, stops)
}

func findBestRouteHelper(curRoute []Stop, stopsLeft []Stop) ([]Stop, time.Duration) {
	if len(stopsLeft) == 0 {
		numberOfRoutes++
		duration := findRouteTime(curRoute)
		fmt.Printf("Routes Tested: %d\nDuration: %v\n\n", numberOfRoutes, duration)
		return curRoute, duration
	}

	var bestRoute []Stop
	var bestDuration time.Duration
	bestDuration = time.Duration(int64(^uint64(0) >> 1))

	for i := range stopsLeft {
		route, duration := findBestRouteHelper(append(curRoute, stopsLeft[i]), removeIndex(i, stopsLeft))
		if duration < bestDuration {
			bestRoute = route
			bestDuration = duration
		}
	}

	return bestRoute, bestDuration
}

func removeIndex(index int, list []Stop) []Stop {
	newList := make([]Stop, 0)
	newList = append(newList, list[:index]...)
	newList = append(newList, list[index+1:]...)
	return newList
}

func findRouteTime(route []Stop) time.Duration {
	var duration time.Duration
	for i := 0; i < len(route)-1; i++ {
		duration += findEdgeTime(route[i], route[i+1], startTimeForRoutes.Add(duration).Unix())
	}
	return duration
}

func findEdgeTime(stopA Stop, stopB Stop, startTime int64) time.Duration {
	req := &maps.DistanceMatrixRequest{
		Origins:       []string{stopA.getCoordinateString()},
		Destinations:  []string{stopB.getCoordinateString()},
		DepartureTime: strconv.FormatInt(startTime, 10),
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

	if resp.Rows[0].Elements[0].Status != "OK" {
		log.Fatalf("Non OK element status: %s\n%s\n", resp.Rows[0].Elements[0].Status, err)
	}

	duration := resp.Rows[0].Elements[0].Duration
	durationInTraffic := resp.Rows[0].Elements[0].DurationInTraffic
	if durationInTraffic > duration {
		duration = durationInTraffic
	}

	return duration
}
