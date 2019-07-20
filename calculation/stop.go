package calculation

import "github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"

// Stop is used to calculate the best route of stops to take
type Stop struct {
	Name           string
	WalkToNextStop bool
}

func dataCollectionStopToCalcStop(dataStops []datacollection.Stop) []Stop {
	stops := make([]Stop, len(dataStops))

	for i, stop := range dataStops {
		stops[i].Name = stop.Name
	}

	return stops
}

func canWalkToNextStop(route []Stop, nextStop Stop, timeFunctions CubicSplineFunctionsHolder, isLastStop bool) bool {
	if len(route) == 0 { // can't go if it is first stop
		return false
	}
	if isLastStop { // can't go if it is last stop
		return false
	}

	lastStop := route[len(route)-1]
	if _, ok := timeFunctions[datacollection.GetEdgeKeyWalking(lastStop.Name, nextStop.Name)]; !ok { // no walking edge for this one
		return false
	}

	if len(route) == 1 { // there is a walking edge and there has only been one stop so far
		return true
	}

	stopBeforeLast := route[len(route)-1]
	if stopBeforeLast.WalkToNextStop { // walked to the last stop, can't walk again
		return false
	}

	return true
}
