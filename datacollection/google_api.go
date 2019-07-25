package datacollection

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"googlemaps.github.io/maps"
)

const (
	// EdgeDataFileDateFormat time format for the date to save to the EdgeData json files
	EdgeDataFileDateFormat string        = "2006-01-02"
	MaxDuration            time.Duration = time.Hour * 1000
)

var allowedTransitLines = map[string]bool{
	"Blue Line":        true,
	"Red Line":         true,
	"Orange Line":      true,
	"Green Line":       true,
	"Green Line B":     true,
	"Green Line C":     true,
	"Green Line D":     true,
	"Green Line E":     true,
	"Mattapan Trolley": true,
}

// EdgeDataTimeLocation location that dates and times should be in
var EdgeDataTimeLocation, _ = time.LoadLocation("America/New_York")

// EdgeTimes stores the time each edge took to traverse (unix time -> duration for that time)
type EdgeTimes map[int64]time.Duration

// Edges a map of all edges by edge key
type Edges map[string]EdgeTimes

// GetEdgeKey returns the map key for an edge between two stops
func GetEdgeKey(stopAName, stopBName string) string {
	return stopAName + ":" + stopBName
}

// GetEdgeKeyWalking returns the map key for an edge between two stops with walking
func GetEdgeKeyWalking(stopAName, stopBName string) string {
	return stopAName + ":" + stopBName + "-Walking"
}

// GetDateFromEdgeDataFilename return the date from the filename of a EdgeData file
func GetDateFromEdgeDataFilename(filename string) (time.Time, error) {
	baseFilename := filepath.Base(filename)
	baseFilename = strings.Replace(baseFilename, ".json", "", -1)
	date := strings.Split(baseFilename, " ")[1]

	// midnight in EST
	return ParseDateToEST(date)
}

// ParseDateToEST takes in a string of format EdgeDataTimeLocation and returns a time.Time variable in EST
func ParseDateToEST(date string) (time.Time, error) {
	t, err := time.Parse(EdgeDataFileDateFormat, date)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(EdgeDataTimeLocation).Add(time.Hour * 4), nil
}

// GetTransitDataFilename returns the filename that represents this data
func GetTransitDataFilename(startTime time.Time, interval time.Duration) string {
	return fmt.Sprintf("datacollection/EdgeData %s.json", startTime.Format(EdgeDataFileDateFormat))
}

// WriteEdgeDataToFile writes the edge data to a file
func WriteEdgeDataToFile(edges Edges, startTime time.Time, interval time.Duration) error {
	data, _ := json.Marshal(edges)
	return ioutil.WriteFile(GetTransitDataFilename(startTime, interval), data, 0644)
}

func readAPIKey(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Can't read api key: %v", err)
	}
	return string(data)
}

// GetTransitDataForAnEdgeWithGoogleAPI returns edge data for a particular edge stopA->stopB
func GetTransitDataForAnEdgeWithGoogleAPI(stopA, stopB *Stop, startTime, endTime time.Time, interval time.Duration, apiKeyFile string) Edges {
	specialEdges, err := ReadSpecialEdgesFromFile(SpecialEdgesFileWithLocationData)
	if err != nil {
		log.Fatalf("Failed to import special edges data: %s", err)
	}

	apiKey := readAPIKey(apiKeyFile)
	mapsClient, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to initialize maps client: %s", err)
	}

	return getTransitDataForAnEdgeWithClient(stopA, stopB, startTime, endTime, interval, mapsClient, specialEdges)
}

func getTransitDataForAnEdgeWithClient(stopA, stopB *Stop, startTime, endTime time.Time, interval time.Duration, mapsClient *maps.Client, specialEdges SpecialEdges) Edges {
	edges := make(Edges)

	if specialEdge, ok := specialEdges[GetEdgeKey(stopA.Name, stopB.Name)]; ok {
		// Walking
		if !specialEdge.NoWalking {
			edges[GetEdgeKeyWalking(stopA.Name, stopB.Name)] = makeEdgeTimeAPICalls(stopA, stopB, interval, startTime, endTime, mapsClient)
		}

		// Transit
		midStop := specialEdge.Stop
		aToMid := makeEdgeTimeAPICalls(stopA, midStop, interval, startTime, endTime, mapsClient)
		midToB := makeEdgeTimeAPICalls(midStop, stopB, interval, startTime, endTime, mapsClient)

		fullEdgeTime := make(EdgeTimes, len(aToMid))
		for k := range aToMid {
			fullEdgeTime[k] = aToMid[k] + midToB[k]
		}
		edges[GetEdgeKey(stopA.Name, stopB.Name)] = fullEdgeTime
	} else {
		edges[GetEdgeKey(stopA.Name, stopB.Name)] = makeEdgeTimeAPICalls(stopA, stopB, interval, startTime, endTime, mapsClient)
	}

	return edges
}

// GetTransitDataWithGoogleAPI generates a json file of distance data of edges
func GetTransitDataWithGoogleAPI(startTime, endTime time.Time, interval time.Duration, apiKeyFile string) {
	stops, err := ImportStopsFromFile(StopLocations)
	if err != nil {
		log.Fatalf("Failed to import stop location data: %s", err)
	}

	specialEdges, err := ReadSpecialEdgesFromFile(SpecialEdgesFileWithLocationData)
	if err != nil {
		log.Fatalf("Failed to import special edges data: %s", err)
	}

	apiKey := readAPIKey(apiKeyFile)
	mapsClient, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to initialize maps client: %s", err)
	}

	numStops := len(stops) - 1
	edges := make(Edges, numStops*numStops)
	for i, stopA := range stops {
		for j, stopB := range stops {
			if i != j {
				newEdges := getTransitDataForAnEdgeWithClient(stopA, stopB, startTime, endTime, interval, mapsClient, specialEdges)
				for k, v := range newEdges {
					edges[k] = v
				}
			}
		}
		log.Printf("Done with %s", stopA.Name)
	}

	err = WriteEdgeDataToFile(edges, startTime, interval)
	if err != nil {
		log.Fatalf("Failed to write edge data file: %s", err)
	}
}

// TODO: Finish
func ReconstructRoute(route []Stop, startTime time.Time, apiKeyFile string) {
	stops, err := ImportStopsFromFile(StopLocations)
	if err != nil {
		log.Fatalf("Failed to import stop location data: %s", err)
	}
	stopMap := make(map[string]*Stop, len(stops))
	for _, stop := range stops {
		stopMap[stop.Name] = stop
	}

	specialEdges, err := ReadSpecialEdgesFromFile(SpecialEdgesFileWithLocationData)
	if err != nil {
		log.Fatalf("Failed to import special edges data: %s", err)
	}

	apiKey := readAPIKey(apiKeyFile)
	mapsClient, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to initialize maps client: %s", err)
	}

	printRoute := func(a, b *Stop, r maps.Route) {
		fmt.Printf("--->%s to %s\n", a.Name, b.Name)
		for _, leg := range r.Legs {
			for _, step := range leg.Steps {
				transitDetails := step.TransitDetails
				if transitDetails == nil {
					for _, substep := range step.Steps {
						if substep.HTMLInstructions != "" {
							fmt.Printf("%s\n", substep.HTMLInstructions)
						}
					}
					continue
				}
				fmt.Printf("depart: %s arrive: %s\n", transitDetails.DepartureStop.Name, transitDetails.ArrivalStop.Name)
			}
		}
	}

	timeToGo := startTime
	for i := 0; i < len(route)-1; i++ {
		stopA := stopMap[route[i].Name]
		stopB := stopMap[route[i+1].Name]

		if specialEdge, ok := specialEdges[GetEdgeKey(stopA.Name, stopB.Name)]; ok && !route[i].WalkToNextStop {
			midStop := specialEdge.Stop

			routeAToMid := getDirectionBetweenEdgesAPICall(stopA, midStop, timeToGo, mapsClient)[0]
			timeToGo = routeAToMid.Legs[0].ArrivalTime
			routeMidToB := getDirectionBetweenEdgesAPICall(midStop, stopB, timeToGo, mapsClient)[0]
			timeToGo = routeMidToB.Legs[0].ArrivalTime
			printRoute(stopA, midStop, routeAToMid)
			printRoute(midStop, stopB, routeMidToB)
		} else {
			routesAToB := getDirectionBetweenEdgesAPICall(stopA, stopB, timeToGo, mapsClient)
			timeToGo = routesAToB[0].Legs[0].ArrivalTime
			printRoute(stopA, stopB, routesAToB[0])
			if stopA.Name == "Forest Hills" {
				fmt.Printf("\n\n\n%f\n\n\n%d\n", findEdgeTime(stopA, stopB, timeToGo, mapsClient).Minutes(), len(routesAToB))
			}
		}

		fmt.Println(timeToGo)
	}
}

// ImportEdgeData gets the edge data from a file
func ImportEdgeData(filename string) (Edges, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	edges := make(Edges)
	err = json.Unmarshal(data, &edges)
	if err != nil {
		return nil, err
	}
	return edges, nil
}

// TODO: Implement
func getDirectionBetweenEdgesAPICall(stopA, stopB *Stop, departTime time.Time, mapsClient *maps.Client) []maps.Route {
	req := &maps.DirectionsRequest{
		Origin:        stopA.LongitudeCommaLatitude,
		Destination:   stopB.LongitudeCommaLatitude,
		DepartureTime: strconv.FormatInt(departTime.Unix(), 10),
		Mode:          maps.TravelModeTransit,
		TransitMode: []maps.TransitMode{
			maps.TransitModeSubway,
			maps.TransitModeTram,
		},
	}

	var routes []maps.Route
	count := 0
	for {
		if count > 5 {
			log.Fatalf("More than 5 retries on query.")
		}

		var err error
		routes, _, err = mapsClient.Directions(context.Background(), req)
		if err != nil {
			log.Printf("directions api error: %s", err)
			count++
			continue
		}

		break
	}

	return routes
}

func makeEdgeTimeAPICalls(stopA, stopB *Stop, interval time.Duration, startTime, endTime time.Time, mapsClient *maps.Client) EdgeTimes {
	edgeTimes := make(EdgeTimes)
	for curTime := startTime; curTime.Before(endTime) || curTime.Equal(endTime); curTime = curTime.Add(interval) {
		edgeTime := findEdgeTime(stopA, stopB, curTime, mapsClient)
		if edgeTime == MaxDuration {
			if !curTime.Equal(startTime) {
				edgeTime = edgeTimes[curTime.Add(-interval).Unix()]
			}
			// if startTime, just put in the MaxDuration time... might not be any trains running in the morning or something
		}
		edgeTimes[curTime.Unix()] = edgeTime
	}
	return edgeTimes
}

func findEdgeTime(stopA, stopB *Stop, startTime time.Time, mapsClient *maps.Client) time.Duration {
	req := &maps.DirectionsRequest{
		Origin:        stopA.LongitudeCommaLatitude,
		Destination:   stopB.LongitudeCommaLatitude,
		DepartureTime: strconv.FormatInt(startTime.Unix(), 10),
		Mode:          maps.TravelModeTransit,
		Alternatives:  true,
		TransitMode: []maps.TransitMode{
			maps.TransitModeSubway,
			maps.TransitModeTram,
			maps.TransitModeRail,
		},
	}

	var routes []maps.Route
	count := 0
	for {
		if count > 5 {
			log.Fatalf("More than 5 retries on query.")
		}

		var err error
		routes, _, err = mapsClient.Directions(context.Background(), req)
		if err != nil {
			log.Printf("directions api error: %s", err)
			count++
			continue
		}

		break
	}

	bestDur := MaxDuration
	for _, route := range routes {
		var routeDur time.Duration
		leg := route.Legs[0] // transit will always only have one leg

		routeDur = leg.ArrivalTime.Sub(startTime)
		for _, step := range leg.Steps {
			if step.TransitDetails != nil && !allowedTransitLines[step.TransitDetails.Line.Name] {
				// must be bus or something like that
				routeDur = MaxDuration
				break
			}
		}
		if routeDur < bestDur {
			bestDur = routeDur
		}
	}

	return bestDur
}
