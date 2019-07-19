package calculation

import (
	"fmt"
	"log"
	"time"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"
)

type Calculator struct {
	startTimeForRoutes time.Time
	numberOfRoutes     int64
	timeFunctions      CubicSplineFunctionsHolder
	startTime          time.Time
	bestTime           time.Duration
	bestRoute          []datacollection.Stop
}

func NewCalculator(timeFunctions CubicSplineFunctionsHolder) (*Calculator, error) {
	return &Calculator{timeFunctions: timeFunctions, bestTime: (time.Hour * 1000)}, nil
}

// FindBestRoute finds the fastest route to traverse every stop, every stop must have an edge to every other stop
func (c *Calculator) FindBestRoute(stops []datacollection.Stop, startTime time.Time) ([]datacollection.Stop, time.Duration) {
	c.numberOfRoutes = 0
	c.startTimeForRoutes = startTime

	route := make([]datacollection.Stop, 0)
	c.startTime = time.Now()
	return c.findBestRouteHelper(route, stops)
}

func printStops(stops []datacollection.Stop) string {
	str := "["
	for _, stop := range stops {
		str += stop.Name + ","
	}
	return str[:len(str)-1] + "]"
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
			elapsed := time.Since(c.startTime)
			c.startTime = time.Now()
			fmt.Printf("Routes Tested: %d\nBest Time: %v Route: %s\n", c.numberOfRoutes, c.bestTime, printStops(c.bestRoute))
			fmt.Printf("Time taken to calculate: %s\n\n", elapsed)
		}
		return curRoute, duration
	}

	var bestRoute []datacollection.Stop
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

func removeIndex(index int, list []datacollection.Stop) []datacollection.Stop {
	newList := make([]datacollection.Stop, 0)
	newList = append(newList, list[:index]...)
	newList = append(newList, list[index+1:]...)
	return newList
}

func (c *Calculator) findRouteTime(route []datacollection.Stop) time.Duration {
	var duration time.Duration
	for i := 0; i < len(route)-1; i++ {
		duration += c.findEdgeTime(route[i], route[i+1], c.startTimeForRoutes.Add(duration))
	}
	return duration
}

func (c *Calculator) findEdgeTime(stopA, stopB datacollection.Stop, startTime time.Time) time.Duration {
	cubicSpline := c.timeFunctions[datacollection.GetEdgeKey(&stopA, &stopB)]
	dur := GetDurationForEdgeFromCubicSpline(cubicSpline, startTime)
	if dur.Hours() > 3 {
		log.Printf("time: %s stopA: %s stopB: %s bad: %d", startTime.String(), stopA.Name, stopB.Name, int(dur.Hours()))
	}
	return dur
}
