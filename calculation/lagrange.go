package calculation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/DzananGanic/numericalgo/interpolate"
	"github.com/DzananGanic/numericalgo/interpolate/lagrange"
)

type LagrangeFunctionsHolder map[string]*lagrange.Lagrange

// MakeLagrangeFunctionForAllEdges returns a map of edges to lagrange time functions, key is the name of both stops seperated by a colon
func (c *Calculator) MakeLagrangeFunctionForAllEdges(stops []Stop, interval time.Duration, startTime, endTime time.Time) LagrangeFunctionsHolder {
	numStops := len(stops) - 1
	lagrangeFunctions := make(LagrangeFunctionsHolder, numStops*numStops)
	for _, stopA := range stops {
		for _, stopB := range stops {
			if stopA != stopB {
				lagrangeFunctions[stopA.Name+":"+stopB.Name] = c.MakeLagrangeFunctionForEdge(stopA, stopB, interval, startTime, endTime)
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

func (c *Calculator) MakeLagrangeFunctionForEdge(stopA, stopB Stop, interval time.Duration, startTime, endTime time.Time) *lagrange.Lagrange {
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

	return approx
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
