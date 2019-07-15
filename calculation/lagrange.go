package calculation

import (
	"fmt"
	"time"

	"github.com/DzananGanic/numericalgo/interpolate"
	"github.com/DzananGanic/numericalgo/interpolate/lagrange"
)

func makeLagrangeFunctionForEdge(stopA, stopB Stop, interval time.Duration, startTime, endTime time.Time) *lagrange.Lagrange {
	x := []float64{}
	y := []float64{}

	for curTime := startTime; curTime.Before(endTime); curTime.Add(interval) {
		newXVal := lagrangeUnitFromTime(curTime)
		x = append(x, newXVal)

		edgeTime := findEdgeTime(stopA, stopB, curTime.Unix())
		newYVal := lagrangeUnitFromDuration(edgeTime)
		y = append(y, newYVal)
	}

	approx := lagrange.New()
	approx.Fit(x, y) // could return error, but only if x and y are different lengths. In this case they won't be

	return approx
}

func getDurationForEdgeFromLagrange(approxFunc *lagrange.Lagrange, startTime time.Time) (time.Duration, error) {
	timeFloat, err := interpolate.WithSingle(approxFunc, lagrangeUnitFromTime(startTime))
	return durationFromLagrangeUnit(timeFloat), err
}

func lagrangeUnitFromTime(t time.Time) float64 {
	return float64((t.Hour() * 60) + (t.Minute()))
}

func lagrangeUnitFromDuration(t time.Duration) float64 {
	return float64((t.Hours() * 60) + (t.Minutes()))
}

func durationFromLagrangeUnit(f float64) time.Duration {
	hours := int(f / 60)
	mins := int(f - float64(hours*60))
	durationString := fmt.Sprintf("%dh%dm", hours, mins)
	duration, _ := time.ParseDuration(durationString)
	return duration
}
