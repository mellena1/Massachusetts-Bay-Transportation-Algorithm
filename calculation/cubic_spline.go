package calculation

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/cnkei/gospline"
	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// cubicSplineFunctionsHolder holds cubic splines given their edge
type cubicSplineFunctionsHolder map[string]gospline.Spline

// makeCubicSplineFunctionForAllEdges returns a map of edges to CubicSpline time functions, key is the name of both stops seperated by a colon
func makeCubicSplineFunctionForAllEdges(edges datacollection.Edges) cubicSplineFunctionsHolder {
	cubicSplineFunctions := make(cubicSplineFunctionsHolder)

	for edge, edgeTimings := range edges {
		cubicSplineFunctions[edge] = makeCubicSplineFunctionForEdge(edgeTimings)
	}

	return cubicSplineFunctions
}

// makeCubicSplineFunctionForEdge takes in edge timing data and creates a cubic spline from it
func makeCubicSplineFunctionForEdge(edgeTiming datacollection.EdgeTimes) gospline.Spline {
	type xy struct {
		x float64
		y float64
	}

	xyPoints := []xy{}
	for unixTime, duration := range edgeTiming {
		xTime := time.Unix(unixTime, 0)

		newPoint := xy{
			x: cubicSplineUnitFromTime(xTime),
			y: cubicSplineUnitFromDuration(duration),
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

// PlotEdge makes a graph of one of the edges
func PlotEdge(edges datacollection.Edges, edgeKey string, outputFile string) error {
	splines := makeCubicSplineFunctionForAllEdges(edges)
	cubicSplineFunc, ok := splines[edgeKey]
	if !ok {
		return errors.New("edge key does not exist")
	}

	p, _ := plot.New()
	p.Title.Text = "Timings for edge " + edgeKey

	loc, _ := time.LoadLocation("America/New_York")

	pts := make(plotter.XYs, 100)
	for i := time.Date(2019, time.July, 18, 6, 0, 0, 0, loc); i.Before(time.Date(2019, time.July, 19, 0, 0, 0, 0, loc)); i = i.Add(time.Minute * 5) {
		var pt plotter.XY
		pt.X = cubicSplineUnitFromTime(i)
		dur := getDurationForEdgeFromCubicSpline(cubicSplineFunc, i)
		pt.Y = cubicSplineUnitFromDuration(dur)
		pts = append(pts, pt)
	}

	plotutil.AddLinePoints(p, "CubicSpline", pts)
	return p.Save(10*vg.Inch, 10*vg.Inch, outputFile)
}

// PlotAllEdges makes a graph with all edges on it
func PlotAllEdges(edges datacollection.Edges, filename string) error {
	cubicSplineFuncs := makeCubicSplineFunctionForAllEdges(edges)

	p, _ := plot.New()
	p.Title.Text = "All Route Times"

	loc, _ := time.LoadLocation("America/New_York")

	for k, v := range cubicSplineFuncs {
		currentCubicSplineFunc := v
		pts := make(plotter.XYs, 100)

		for i := time.Date(2019, time.July, 18, 6, 0, 0, 0, loc); i.Before(time.Date(2019, time.July, 19, 0, 0, 0, 0, loc)); i = i.Add(time.Minute * 5) {
			var pt plotter.XY
			pt.X = cubicSplineUnitFromTime(i)
			dur := getDurationForEdgeFromCubicSpline(currentCubicSplineFunc, i)
			pt.Y = cubicSplineUnitFromDuration(dur)
			pts = append(pts, pt)
		}
		plotutil.AddLinePoints(p, k, pts)
	}

	return p.Save(24*vg.Inch, 24*vg.Inch, filename)
}

func getDurationForEdgeFromCubicSpline(approxFunc gospline.Spline, startTime time.Time) time.Duration {
	timeFloat := approxFunc.At(cubicSplineUnitFromTime(startTime))
	return durationFromCubicSplineUnit(timeFloat)
}

func cubicSplineUnitFromTime(t time.Time) float64 {
	hour := t.Hour()
	if hour <= 4 {
		hour += 24
	}
	return float64((hour * 60) + t.Minute())
}

func cubicSplineUnitFromDuration(t time.Duration) float64 {
	return t.Minutes()
}

func durationFromCubicSplineUnit(f float64) time.Duration {
	hours := int(f / 60)
	mins := int(f - float64(hours*60))
	durationString := fmt.Sprintf("%dh%dm", hours, mins)
	duration, _ := time.ParseDuration(durationString)
	return duration
}
