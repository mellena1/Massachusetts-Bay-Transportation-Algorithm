package calculation

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/cnkei/gospline"
	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type CubicSplineFunctionsHolder map[string]gospline.Spline

// MakeCubicSplineFunctionForAllEdges returns a map of edges to CubicSpline time functions, key is the name of both stops seperated by a colon
func MakeCubicSplineFunctionForAllEdges(edges datacollection.Edges) CubicSplineFunctionsHolder {
	cubicSplineFunctions := make(CubicSplineFunctionsHolder)

	for edge, edgeTimings := range edges {
		cubicSplineFunctions[edge] = MakeCubicSplineFunctionForEdge(edgeTimings)
	}

	return cubicSplineFunctions
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

func MakeCubicSplineFunctionForEdge(edgeTiming datacollection.EdgeTimes) gospline.Spline {
	type xy struct {
		x float64
		y float64
	}

	xyPoints := []xy{}
	for unixTime, duration := range edgeTiming {
		xTime := time.Unix(unixTime, 0)

		newPoint := xy{
			x: CubicSplineUnitFromTime(xTime),
			y: CubicSplineUnitFromDuration(duration),
		}
		xyPoints = append(xyPoints, newPoint)
	}
	// x values must be in ascending order
	sort.Slice(xyPoints, func(i int, j int) bool {
		return xyPoints[i].x < xyPoints[j].x
	})

	x := make([]float64, len(xyPoints))
	y := make([]float64, len(xyPoints))
	for i, point := range xyPoints {
		x[i] = point.x
		y[i] = point.y
	}

	return gospline.NewCubicSpline(x, y)
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
