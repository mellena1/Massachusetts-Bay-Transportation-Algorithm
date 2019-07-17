package calculation

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/cnkei/gospline"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"googlemaps.github.io/maps"
)

type CubicSplineFunctionsHolder map[string]gospline.Spline

func getLFHKey(stopA, stopB Stop) string {
	return stopA.Name + ":" + stopB.Name
}

type CubicSpline struct {
	mapsClient *maps.Client
}

func NewCubicSpline(apiKey string) (*CubicSpline, error) {
	mapsClient, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &CubicSpline{mapsClient: mapsClient}, nil
}

type EdgeTimes map[int64]time.Duration
type Edges map[string]EdgeTimes

// MakeCubicSplineFunctionForAllEdges returns a map of edges to CubicSpline time functions, key is the name of both stops seperated by a colon
func (l *CubicSpline) MakeCubicSplineFunctionForAllEdges(stops []Stop, interval time.Duration, startTime, endTime time.Time, edges Edges) CubicSplineFunctionsHolder {
	numStops := len(stops) - 1
	cubicSplineFunctions := make(CubicSplineFunctionsHolder, numStops*numStops)
	for i, stopA := range stops {
		for j, stopB := range stops {
			if i != j {
				cubicSplineFunctions[getLFHKey(stopA, stopB)] = l.MakeCubicSplineFunctionForEdge(stopA, stopB, interval, startTime, endTime, edges)
			}
		}
		log.Printf("Done with %s", stopA.Name)
	}

	return cubicSplineFunctions
}

func (l *CubicSpline) SaveAPICalls(stops []Stop, interval time.Duration, startTime, endTime time.Time, filename string) {
	numStops := len(stops) - 1
	edges := make(Edges, numStops*numStops)
	for i, stopA := range stops {
		for j, stopB := range stops {
			if i != j {
				edges[getLFHKey(stopA, stopB)] = l.makeAPICall(stopA, stopB, interval, startTime, endTime)
			}
		}
		log.Printf("Done with %s", stopA.Name)
	}

	data, _ := json.Marshal(edges)
	ioutil.WriteFile(filename, data, 0644)
}

func ReadAPICalls(filename string) (Edges, error) {
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

func (l *CubicSpline) makeAPICall(stopA, stopB Stop, interval time.Duration, startTime, endTime time.Time) EdgeTimes {
	edgeTimes := make(EdgeTimes)
	for curTime := startTime; curTime.Before(endTime) || curTime.Equal(endTime); curTime = curTime.Add(interval) {
		unixTime := curTime.Unix()
		edgeTimes[unixTime] = l.findEdgeTime(stopA, stopB, unixTime)
	}
	return edgeTimes
}

func WriteCubicSplineFunctionsToFile(cubicSplines CubicSplineFunctionsHolder, filename string) error {
	data, err := json.Marshal(cubicSplines)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	return err
}

func ReadCubicSplineFunctionsFromFile(filename string) (CubicSplineFunctionsHolder, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cubicSplines := make(CubicSplineFunctionsHolder)
	err = json.Unmarshal(data, &cubicSplines)
	return cubicSplines, err
}

func PlotCubicSplineFunc(cubicSplineFunc gospline.Spline, filename string) error {
	p, err := plot.New()
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("America/New_York")

	pts := make(plotter.XYs, 100)
	for i := time.Date(2019, time.July, 18, 6, 0, 0, 0, loc); i.Before(time.Date(2019, time.July, 19, 0, 0, 0, 0, loc)); i = i.Add(time.Minute * 5) {
		var pt plotter.XY
		pt.X = CubicSplineUnitFromTime(i)
		dur := GetDurationForEdgeFromCubicSpline(cubicSplineFunc, i)
		pt.Y = CubicSplineUnitFromDuration(dur)
		pts = append(pts, pt)
	}

	log.Print("all points made")
	plotutil.AddLinePoints(p, "CubicSpline", pts)

	return p.Save(10*vg.Inch, 10*vg.Inch, filename)
}

func (l *CubicSpline) MakeCubicSplineFunctionForEdge(stopA, stopB Stop, interval time.Duration, startTime, endTime time.Time, edges Edges) gospline.Spline {
	x := []float64{}
	y := []float64{}

	for curTime := startTime; curTime.Before(endTime) || curTime.Equal(endTime); curTime = curTime.Add(interval) {
		newXVal := CubicSplineUnitFromTime(curTime)
		x = append(x, newXVal)

		edgeTime := edges[getLFHKey(stopA, stopB)][curTime.Unix()]
		newYVal := CubicSplineUnitFromDuration(edgeTime)
		y = append(y, newYVal)
	}

	approx := gospline.NewCubicSpline(x, y)
	// approx.Fit(x, y) // could return error, but only if x and y are different lengths. In this case they won't be

	return approx
}

func GetDurationForEdgeFromCubicSpline(approxFunc gospline.Spline, startTime time.Time) time.Duration {
	timeFloat := approxFunc.At(CubicSplineUnitFromTime(startTime))
	return durationFromCubicSplineUnit(timeFloat)
}

func CubicSplineUnitFromTime(t time.Time) float64 {
	hour := t.Hour()
	if hour <= 4 {
		hour += 24
	}
	return float64((hour * 60) + t.Minute())
}

func CubicSplineUnitFromDuration(t time.Duration) float64 {
	return t.Minutes()
}

func durationFromCubicSplineUnit(f float64) time.Duration {
	hours := int(f / 60)
	mins := int(f - float64(hours*60))
	durationString := fmt.Sprintf("%dh%dm", hours, mins)
	duration, _ := time.ParseDuration(durationString)
	return duration
}

func (l *CubicSpline) findEdgeTime(stopA Stop, stopB Stop, startTime int64) time.Duration {
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
