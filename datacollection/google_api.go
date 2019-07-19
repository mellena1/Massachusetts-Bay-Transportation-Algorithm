package datacollection

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"googlemaps.github.io/maps"
)

// EdgeTimes stores the time each edge took to traverse
type EdgeTimes map[int64]time.Duration

// Edges a map of all edges by edge key
type Edges map[string]EdgeTimes

// GetEdgeKey returns the map key for an edge between two stops
func GetEdgeKey(stopA, stopB *Stop) string {
	return stopA.Name + ":" + stopB.Name
}

// GetTransitDataFilename returns the filename that represents this data
func GetTransitDataFilename(startTime time.Time, interval time.Duration) string {
	return "EdgeData StartTime:" + strconv.FormatInt(startTime.Unix(), 10) + " Interval:" + interval.String()
}

func readAPIKey() string {
	data, err := ioutil.ReadFile("apikey.secret")
	if err != nil {
		log.Fatalf("Can't read api key: %v", err)
	}
	return string(data)
}

// GetTransitDataWithGoogleAPI generates a json file of distance data of edges
func GetTransitDataWithGoogleAPI(startTime, endTime time.Time, interval time.Duration) {
	stops, err := ImportStopsFromFile(StopLocations)
	if err != nil {
		log.Fatalf("Failed to import stop location data: %s", err)
	}

	apiKey := readAPIKey()
	mapsClient, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to initialize maps client: %s", err)
	}

	numStops := len(stops) - 1
	edges := make(Edges, numStops*numStops)
	for i, stopA := range stops {
		for j, stopB := range stops {
			if i != j {
				edges[GetEdgeKey(stopA, stopB)] = makeAPICall(stopA, stopB, interval, startTime, endTime, mapsClient)
			}
		}
		log.Printf("Done with %s", stopA.Name)
	}

	data, _ := json.Marshal(edges)
	ioutil.WriteFile(GetTransitDataFilename(startTime, interval), data, 0644)
}

func makeAPICall(stopA, stopB *Stop, interval time.Duration, startTime, endTime time.Time, mapsClient *maps.Client) EdgeTimes {
	edgeTimes := make(EdgeTimes)
	for curTime := startTime; curTime.Before(endTime) || curTime.Equal(endTime); curTime = curTime.Add(interval) {
		unixTime := curTime.Unix()
		edgeTimes[unixTime] = findEdgeTime(stopA, stopB, unixTime, mapsClient)
	}
	return edgeTimes
}

func findEdgeTime(stopA, stopB *Stop, startTime int64, mapsClient *maps.Client) time.Duration {
	req := &maps.DistanceMatrixRequest{
		Origins:       []string{stopA.LongitudeCommaLatitude},
		Destinations:  []string{stopB.LongitudeCommaLatitude},
		DepartureTime: strconv.FormatInt(startTime, 10),
		Mode:          maps.TravelModeTransit,
		TransitMode: []maps.TransitMode{
			maps.TransitModeSubway,
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
		resp, err = mapsClient.DistanceMatrix(context.Background(), req)
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