package calculation

import (
	"fmt"
	"time"

	"github.com/cnkei/gospline"
	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"
)

// Calculator contains data for route calculation
type Calculator struct {
	startTimeForRoutes time.Time
	numberOfRoutes     int64
	timeFunctions      CubicSplineFunctionsHolder
	startTime          time.Time
	bestTime           time.Duration
	bestRoute          []Stop
}

// NewCalculator returns a new calculator object
func NewCalculator(edgeData datacollection.Edges) (*Calculator, error) {
	return &Calculator{timeFunctions: MakeCubicSplineFunctionForAllEdges(edgeData), bestTime: (time.Hour * 1000)}, nil
}

// FindBestRoute finds the fastest route to traverse every stop, every stop must have an edge to every other stop
func (c *Calculator) FindBestRoute(stops []datacollection.Stop, startTime time.Time) ([]Stop, time.Duration) {
	convStops := dataCollectionStopToCalcStop(stops)

	c.numberOfRoutes = 0
	c.startTimeForRoutes = startTime

	route := make([]Stop, 0)
	c.startTime = time.Now()
	return c.findBestRouteHelper(route, convStops)
}

func PrintStops(stops []Stop) string {
	str := "["
	for _, stop := range stops {
		if stop.WalkToNextStop {
			str += stop.Name + "-walk,"
		} else {
			str += stop.Name + ","
		}
	}
	return str[:len(str)-1] + "]"
}

func (c *Calculator) findBestRouteHelper(curRoute, stopsLeft []Stop) ([]Stop, time.Duration) {
	if len(stopsLeft) == 0 {
		c.numberOfRoutes++
		duration := c.findRouteTime(curRoute)
		if duration.Minutes() < c.bestTime.Minutes() {
			c.bestTime = duration
			c.bestRoute = curRoute
		}
		if c.numberOfRoutes%1000000 == 0 {
			elapsed := time.Since(c.startTime)
			c.startTime = time.Now()
			fmt.Printf("Routes Tested: %d\nBest Time: %v Route: %s\n", c.numberOfRoutes, c.bestTime, PrintStops(c.bestRoute))
			fmt.Printf("Time taken to calculate: %s\n\n", elapsed)
		}
		return curRoute, duration
	}

	var bestRoute []Stop
	var bestDuration time.Duration
	bestDuration = time.Duration(int64(^uint64(0) >> 1))

	for i := range stopsLeft {
		// newRoute := cloneRouteSlice(curRoute)
		route, duration := c.findBestRouteHelper(append(curRoute, stopsLeft[i]), removeIndex(i, stopsLeft))
		if duration < bestDuration {
			bestRoute = route
			bestDuration = duration
		}
		if canWalkToNextStop(curRoute, stopsLeft[i], c.timeFunctions, len(stopsLeft) == 1) {
			newRoute := append(curRoute, stopsLeft[i])
			newRoute[len(newRoute)-2].WalkToNextStop = true
			route, duration := c.findBestRouteHelper(newRoute, removeIndex(i, stopsLeft))
			if duration < bestDuration {
				bestRoute = route
				bestDuration = duration
			}
		}
	}

	return bestRoute, bestDuration
}

func cloneRouteSlice(route []Stop) []Stop {
	newRoute := make([]Stop, len(route))
	copy(newRoute, route)
	return newRoute
}

func removeIndex(index int, list []Stop) []Stop {
	newList := make([]Stop, 0)
	newList = append(newList, list[:index]...)
	newList = append(newList, list[index+1:]...)
	return newList
}

func (c *Calculator) findRouteTime(route []Stop) time.Duration {
	var duration time.Duration
	for i := 0; i < len(route)-1; i++ {
		duration += c.findEdgeTime(route[i], route[i+1], c.startTimeForRoutes.Add(duration))
	}
	return duration
}

func (c *Calculator) findEdgeTime(stopA, stopB Stop, startTime time.Time) time.Duration {
	var cubicSpline gospline.Spline
	if stopA.WalkToNextStop {
		cubicSpline = c.timeFunctions[datacollection.GetEdgeKeyWalking(stopA.Name, stopB.Name)]
	} else {
		cubicSpline = c.timeFunctions[datacollection.GetEdgeKey(stopA.Name, stopB.Name)]
	}
	dur := GetDurationForEdgeFromCubicSpline(cubicSpline, startTime)
	return dur
}
