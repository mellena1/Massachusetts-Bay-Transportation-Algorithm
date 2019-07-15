package calculation

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"googlemaps.github.io/maps"
)

type Calculator struct {
	mapsClient         *maps.Client
	startTimeForRoutes time.Time
	numberOfRoutes     int64
}

func NewCalculator(apiKey string) (*Calculator, error) {
	mapsClient, err := maps.NewClient(maps.WithAPIKey(apiKey))
	return &Calculator{mapsClient: mapsClient}, err
}

// FindBestRoute finds the fastest route to traverse every stop, every stop must have an edge to every other stop
func (c *Calculator) FindBestRoute(stops []Stop, startTime time.Time) ([]Stop, time.Duration) {
	c.numberOfRoutes = 0
	c.startTimeForRoutes = startTime

	route := make([]Stop, 0)
	return c.findBestRouteHelper(route, stops)
}

func (c *Calculator) findBestRouteHelper(curRoute []Stop, stopsLeft []Stop) ([]Stop, time.Duration) {
	if len(stopsLeft) == 0 {
		c.numberOfRoutes++
		duration := c.findRouteTime(curRoute)
		fmt.Printf("Routes Tested: %d\nDuration: %v\n\n", c.numberOfRoutes, duration)
		return curRoute, duration
	}

	var bestRoute []Stop
	var bestDuration time.Duration
	bestDuration = time.Duration(int64(^uint64(0) >> 1))

	for i := range stopsLeft {
		route, duration := c.findBestRouteHelper(append(curRoute, stopsLeft[i]), removeIndex(i, stopsLeft))
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

func (c *Calculator) findRouteTime(route []Stop) time.Duration {
	var duration time.Duration
	for i := 0; i < len(route)-1; i++ {
		duration += c.findEdgeTime(route[i], route[i+1], c.startTimeForRoutes.Add(duration).Unix())
	}
	return duration
}

func (c *Calculator) findEdgeTime(stopA Stop, stopB Stop, startTime int64) time.Duration {
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

	var resp *maps.DistanceMatrixResponse
	count := 0
	for {
		if count > 5 {
			log.Fatalf("More than 5 retries on query.")
		}

		var err error
		resp, err = c.mapsClient.DistanceMatrix(context.Background(), req)
		if err != nil {
			log.Fatalf("fatal error: %s", err)
		}

		if resp.Rows[0].Elements[0].Status != "OK" {
			fmt.Printf("Elements Status: %v\n\n", resp.Rows[0].Elements[0].Status)
			count++
			continue
		}

		break
	}

	duration := resp.Rows[0].Elements[0].Duration
	durationInTraffic := resp.Rows[0].Elements[0].DurationInTraffic
	if durationInTraffic > duration {
		duration = durationInTraffic
	}

	return duration
}
