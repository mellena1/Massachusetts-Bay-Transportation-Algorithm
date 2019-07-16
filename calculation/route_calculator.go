package calculation

import (
	"fmt"
	"log"
	"time"
)

type Calculator struct {
	startTimeForRoutes time.Time
	numberOfRoutes     int64
	timeFunctions      LagrangeFunctionsHolder
	bestTime           time.Duration
}

func NewCalculator(timeFunctions LagrangeFunctionsHolder) (*Calculator, error) {
	return &Calculator{timeFunctions: timeFunctions, bestTime: (time.Hour * 1000)}, nil
}

// FindBestRoute finds the fastest route to traverse every stop, every stop must have an edge to every other stop
func (c *Calculator) FindBestRoute(stops []Stop, startTime time.Time) ([]Stop, time.Duration) {
	c.numberOfRoutes = 0
	c.startTimeForRoutes = startTime

	route := make([]Stop, 0)
	return c.findBestRouteHelper(route, stops)
}

func (c *Calculator) findBestRouteHelper(curRoute []Stop, stopsLeft []Stop) ([]Stop, time.Duration) {
	if len(stopsLeft) == 0 {
		c.numberOfRoutes++
		duration := c.findRouteTime(curRoute)
		if duration.Minutes() < c.bestTime.Minutes() {
			c.bestTime = duration
		}
		if c.numberOfRoutes%1000000 == 0 {
			fmt.Printf("Routes Tested: %d\nBest Time: %v\n\n", c.numberOfRoutes, c.bestTime)
		}
		return curRoute, duration
	}

	var bestRoute []Stop
	var bestDuration time.Duration
	bestDuration = time.Duration(int64(^uint64(0) >> 1))

	for i := range stopsLeft {
		route, duration := c.findBestRouteHelper(append(curRoute, stopsLeft[i]), removeIndex(i, stopsLeft))
		if duration < bestDuration {
			bestRoute = route
			bestDuration = duration
		}
	}

	return bestRoute, bestDuration
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

func (c *Calculator) findEdgeTime(stopA Stop, stopB Stop, startTime time.Time) time.Duration {
	lagrange := c.timeFunctions[getLFHKey(stopA, stopB)]
	dur := GetDurationForEdgeFromLagrange(lagrange, startTime)
	if dur.Hours() > 3 {
		log.Printf("time: %s stopA: %s stopB: %s bad: %d", startTime.String(), stopA.Name, stopB.Name, int(dur.Hours()))
	}
	return dur
}
