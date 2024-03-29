package calculation

import (
	"fmt"
	"time"

	"github.com/cnkei/gospline"
	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"
)

// Calculator contains data for route calculation
type Calculator struct {
	startTimeForRoutes time.Time                  // start all routes at this time
	numberOfRoutes     int64                      // number of routes calculated so far
	timeFunctions      cubicSplineFunctionsHolder // cubic splines
	latestTime         time.Time                  // don't calculate routes past this time
	timer              time.Time                  // timer to see how long each batch of routes takes
	bestTime           time.Duration              // holds the current best time so it can print it out during iteration
	bestRoute          []datacollection.Stop      // holds the current best route so it can print it out during iteration
}

// NewCalculator returns a new calculator object
func NewCalculator(edgeData datacollection.Edges, latestTime time.Time) (*Calculator, error) {
	return &Calculator{timeFunctions: makeCubicSplineFunctionForAllEdges(edgeData), latestTime: latestTime, bestTime: (time.Hour * 1000)}, nil
}

// FindBestRoute finds the fastest route to traverse every stop, every stop must have an edge to every other stop
func (c *Calculator) FindBestRoute(stops []datacollection.Stop, startTime time.Time) ([]datacollection.Stop, time.Duration) {
	c.numberOfRoutes = 0
	c.startTimeForRoutes = startTime

	route := make([]datacollection.Stop, 0)
	c.timer = time.Now()
	return c.findBestRouteHelper(route, stops)
}

// PrintStops prints out an array of the stops in a route
func PrintStops(stops []datacollection.Stop) string {
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

// canWalkToNextStop holds the logic that we decided on for whether or not a route can walk to the next stop it tries to go to
func canWalkToNextStop(route []datacollection.Stop, nextStop datacollection.Stop, timeFunctions cubicSplineFunctionsHolder, isLastStop bool) bool {
	if len(route) <= 1 { // can't go if it is first stop, and don't run if route is len 0
		return false
	}
	if isLastStop { // can't go if it is last stop
		return false
	}

	lastStop := route[len(route)-1]
	if _, ok := timeFunctions[datacollection.GetEdgeKeyWalking(lastStop.Name, nextStop.Name)]; !ok { // no walking edge for this one
		return false
	}

	stopBeforeLast := route[len(route)-2]
	if stopBeforeLast.WalkToNextStop { // walked to the last stop, can't walk again
		return false
	}

	return true
}

func (c *Calculator) findBestRouteHelper(curRoute, stopsLeft []datacollection.Stop) ([]datacollection.Stop, time.Duration) {
	if len(stopsLeft) == 0 {
		c.numberOfRoutes++
		duration := c.findRouteTime(curRoute)
		if duration.Minutes() < c.bestTime.Minutes() {
			c.bestTime = duration
			c.bestRoute = curRoute
		}
		if c.numberOfRoutes%1000000 == 0 {
			elapsed := time.Since(c.timer)
			c.timer = time.Now()
			fmt.Printf("%s - Routes Tested: %d\nBest Time: %v Route: %s\n", c.startTimeForRoutes.Format(time.Kitchen), c.numberOfRoutes, c.bestTime, PrintStops(c.bestRoute))
			fmt.Printf("Time taken to calculate: %s\n\n", elapsed)
		}
		return curRoute, duration
	}

	var bestRoute []datacollection.Stop
	var bestDuration time.Duration
	bestDuration = time.Duration(int64(^uint64(0) >> 1))

	for i := range stopsLeft {
		newRoute := cloneRouteSlice(curRoute)
		route, duration := c.findBestRouteHelper(append(newRoute, stopsLeft[i]), removeIndex(i, stopsLeft))
		if duration < bestDuration {
			bestRoute = route
			bestDuration = duration
		}
		if canWalkToNextStop(curRoute, stopsLeft[i], c.timeFunctions, len(stopsLeft) == 1) {
			newRoute := cloneRouteSlice(curRoute)
			newRoute[len(newRoute)-1].WalkToNextStop = true
			route, duration := c.findBestRouteHelper(append(newRoute, stopsLeft[i]), removeIndex(i, stopsLeft))
			if duration < bestDuration {
				bestRoute = route
				bestDuration = duration
			}
		}
	}

	return bestRoute, bestDuration
}

func cloneRouteSlice(route []datacollection.Stop) []datacollection.Stop {
	newRoute := make([]datacollection.Stop, len(route))
	copy(newRoute, route)
	return newRoute
}

func removeIndex(index int, list []datacollection.Stop) []datacollection.Stop {
	newList := make([]datacollection.Stop, len(list)-1)
	copy(newList, list[:index])
	copy(newList[index:], list[index+1:])
	return newList
}

func (c *Calculator) findRouteTime(route []datacollection.Stop) time.Duration {
	var duration time.Duration
	for i := 0; i < len(route)-1; i++ {
		edgeStartTime := c.startTimeForRoutes.Add(duration)
		if edgeStartTime.After(c.latestTime) {
			return datacollection.MaxDuration
		}
		duration += c.findEdgeTime(route[i], route[i+1], edgeStartTime)
	}
	return duration
}

func (c *Calculator) findEdgeTime(stopA, stopB datacollection.Stop, startTime time.Time) time.Duration {
	var cubicSpline gospline.Spline
	if stopA.WalkToNextStop {
		cubicSpline = c.timeFunctions[datacollection.GetEdgeKeyWalking(stopA.Name, stopB.Name)]
	} else {
		cubicSpline = c.timeFunctions[datacollection.GetEdgeKey(stopA.Name, stopB.Name)]
	}
	dur := getDurationForEdgeFromCubicSpline(cubicSpline, startTime)
	return dur
}
