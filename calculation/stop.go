package calculation

import (
	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"
)

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
