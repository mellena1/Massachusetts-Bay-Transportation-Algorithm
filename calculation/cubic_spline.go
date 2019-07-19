package calculation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/cnkei/gospline"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type CubicSplineFunctionsHolder map[string]gospline.Spline

// GetEdgeKey returns the map key for an edge between two stops
func GetEdgeKey(stopA, stopB Stop) string {
	return stopA.Name + ":" + stopB.Name
}

// MakeCubicSplineFunctionForAllEdges returns a map of edges to CubicSpline time functions, key is the name of both stops seperated by a colon
func MakeCubicSplineFunctionForAllEdges(stops []Stop, interval time.Duration, startTime, endTime time.Time, edges Edges) CubicSplineFunctionsHolder {
	numStops := len(stops) - 1
	cubicSplineFunctions := make(CubicSplineFunctionsHolder, numStops*numStops)
	for i, stopA := range stops {
		for j, stopB := range stops {
			if i != j {
				cubicSplineFunctions[getLFHKey(stopA, stopB)] = MakeCubicSplineFunctionForEdge(stopA, stopB, interval, startTime, endTime, edges)
			}
		}
		log.Printf("Done with %s", stopA.Name)
	}

	return cubicSplineFunctions
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

	plotutil.AddLinePoints(p, "CubicSpline", pts)
	return p.Save(10*vg.Inch, 10*vg.Inch, filename)
}

func PlotAllCubicSplineFuncs(cubicSplineFuncs CubicSplineFunctionsHolder, filename string) error {
	p, err := plot.New()
	p.Title.Text = "All Route Times"

	if err != nil {
		return err
	}
	loc, _ := time.LoadLocation("America/New_York")

	for k, v := range cubicSplineFuncs {
		currentCubicSplineFunc := v
		pts := make(plotter.XYs, 100)

		for i := time.Date(2019, time.July, 18, 6, 0, 0, 0, loc); i.Before(time.Date(2019, time.July, 19, 0, 0, 0, 0, loc)); i = i.Add(time.Minute * 5) {
			var pt plotter.XY
			pt.X = CubicSplineUnitFromTime(i)
			dur := GetDurationForEdgeFromCubicSpline(currentCubicSplineFunc, i)
			pt.Y = CubicSplineUnitFromDuration(dur)
			log.Printf("%f, %f", pt.X, pt.Y)
			pts = append(pts, pt)
		}
		plotutil.AddLinePoints(p, k, pts)
	}

	return p.Save(24*vg.Inch, 24*vg.Inch, filename)
}

func MakeCubicSplineFunctionForEdge(stopA, stopB Stop, interval time.Duration, startTime, endTime time.Time, edges Edges) gospline.Spline {
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
