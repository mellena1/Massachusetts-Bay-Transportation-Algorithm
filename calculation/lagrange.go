package calculation

import (
	"fmt"
	"time"

	"github.com/DzananGanic/numericalgo/interpolate"
	"github.com/DzananGanic/numericalgo/interpolate/lagrange"
)

// MakeLagrangeFunctionForAllEdges returns a map of edges to lagrange time functions, key is the name of both stops seperated by a colon
func MakeLagrangeFunctionForAllEdges(stops []Stop, interval time.Duration, startTime, endTime time.Time) {
	numStops := len(stops) - 1
	lagrangeFunctions := make(map[string]*lagrange.Lagrange, numStops*numStops)
	for _, stopA := range stops {
		for _, stopB := range stops {
			if stopA != stopB {
				lagrangeFunctions[stopA.Name+":"+stopB.Name] = makeLagrangeFunctionForEdge(stopA, stopB, interval, startTime, endTime)
			}
		}
	}
}

type LagrangeApproxEdge struct {
	Lagrange     *lagrange.Lagrange
	StartingStop *Stop
	EndingStop   *Stop
}

func (c *Calculator) MakeLagrangeFunctionForEdge(stopA, stopB Stop, interval time.Duration, startTime, endTime time.Time) *LagrangeApproxEdge {
	x := []float64{}
	y := []float64{}

	for curTime := startTime; curTime.Before(endTime); curTime = curTime.Add(interval) {
		newXVal := lagrangeUnitFromTime(curTime)
		x = append(x, newXVal)

		edgeTime := c.findEdgeTime(stopA, stopB, curTime.Unix())
		newYVal := lagrangeUnitFromDuration(edgeTime)
		y = append(y, newYVal)
	}

	approx := lagrange.New()
	approx.Fit(x, y) // could return error, but only if x and y are different lengths. In this case they won't be

	return &LagrangeApproxEdge{
		Lagrange:     approx,
		StartingStop: &stopA,
		EndingStop:   &stopB,
	}
}

func GetDurationForEdgeFromLagrange(approxFunc *lagrange.Lagrange, startTime time.Time) (time.Duration, error) {
	timeFloat, err := interpolate.WithSingle(approxFunc, lagrangeUnitFromTime(startTime))
	return durationFromLagrangeUnit(timeFloat), err
}

func lagrangeUnitFromTime(t time.Time) float64 {
	return float64((t.Hour() * 60) + (t.Minute()))
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
