package calculation

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/DzananGanic/numericalgo/interpolate"
	"github.com/DzananGanic/numericalgo/interpolate/lagrange"
	"googlemaps.github.io/maps"
)

type LagrangeFunctionsHolder map[string]*lagrange.Lagrange

type Lagrange struct {
	mapsClient *maps.Client
}

func NewLagrange(apiKey string) (*Lagrange, error) {
	mapsClient, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &Lagrange{mapsClient: mapsClient}, nil
}

// MakeLagrangeFunctionForAllEdges returns a map of edges to lagrange time functions, key is the name of both stops seperated by a colon
func (l *Lagrange) MakeLagrangeFunctionForAllEdges(stops []Stop, interval time.Duration, startTime, endTime time.Time) LagrangeFunctionsHolder {
	numStops := len(stops) - 1
	lagrangeFunctions := make(LagrangeFunctionsHolder, numStops*numStops)
	for i, stopA := range stops {
		for j, stopB := range stops {
			if i != j {
				lagrangeFunctions[stopA.Name+":"+stopB.Name] = l.MakeLagrangeFunctionForEdge(stopA, stopB, interval, startTime, endTime)
			}
		}
	}

	return lagrangeFunctions
}

func WriteLangrageFunctionsToFile(lagranges LagrangeFunctionsHolder, filename string) error {
	data, err := json.Marshal(lagranges)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	return err
}

func ReadLagrangeFunctionsFromFile(filename string) (LagrangeFunctionsHolder, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lagranges := make(LagrangeFunctionsHolder)
	err = json.Unmarshal(data, &lagranges)
	return lagranges, err
}

func (l *Lagrange) MakeLagrangeFunctionForEdge(stopA, stopB Stop, interval time.Duration, startTime, endTime time.Time) *lagrange.Lagrange {
	x := []float64{}
	y := []float64{}

	for curTime := startTime; curTime.Before(endTime) || curTime.Equal(endTime); curTime = curTime.Add(interval) {
		newXVal := lagrangeUnitFromTime(curTime)
		x = append(x, newXVal)

		edgeTime := l.findEdgeTime(stopA, stopB, curTime.Unix())
		newYVal := lagrangeUnitFromDuration(edgeTime)
		y = append(y, newYVal)
	}

	approx := lagrange.New()
	approx.Fit(x, y) // could return error, but only if x and y are different lengths. In this case they won't be

	return approx
}

func GetDurationForEdgeFromLagrange(approxFunc *lagrange.Lagrange, startTime time.Time) (time.Duration, error) {
	timeFloat, err := interpolate.WithSingle(approxFunc, lagrangeUnitFromTime(startTime))
	return durationFromLagrangeUnit(timeFloat), err
}

func lagrangeUnitFromTime(t time.Time) float64 {
	hour := t.Hour()
	if hour <= 4 {
		hour += 24
	}
	return float64((hour * 60) + t.Minute())
}

func lagrangeUnitFromDuration(t time.Duration) float64 {
	return t.Minutes()
}

func durationFromLagrangeUnit(f float64) time.Duration {
	hours := int(f / 60)
	mins := int(f - float64(hours*60))
	durationString := fmt.Sprintf("%dh%dm", hours, mins)
	duration, _ := time.ParseDuration(durationString)
	return duration
}

func (l *Lagrange) findEdgeTime(stopA Stop, stopB Stop, startTime int64) time.Duration {
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
		resp, err = l.mapsClient.DistanceMatrix(context.Background(), req)
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
