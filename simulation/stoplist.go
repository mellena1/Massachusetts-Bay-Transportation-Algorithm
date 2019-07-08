package simulation

import (
	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/graph"
)

// StopList holds info about which Stops have been visited
type StopList struct {
	AccessedStops map[*graph.Stop]bool
	Path          []*graph.Stop
	allStops      []*graph.Stop
}

// NewStopList makes a new StopList given all stops
func NewStopList(allStops []*graph.Stop) *StopList {
	return &StopList{
		AccessedStops: make(map[*graph.Stop]bool),
		Path:          []*graph.Stop{},
		allStops:      allStops,
	}
}

func CloneStopList(sl *StopList) *StopList {
	accessedStops := make(map[*graph.Stop]bool, len(sl.AccessedStops))
	for k, v := range sl.AccessedStops {
		accessedStops[k] = v
	}

	path := make([]*graph.Stop, len(sl.Path))
	copy(path, sl.Path)

	return &StopList{
		AccessedStops: accessedStops,
		Path:          path,
		allStops:      sl.allStops,
	}
}

// ArriveAtStop log that a stop has been arrived at
func (sl StopList) ArriveAtStop(stop *graph.Stop) {
	if _, exists := sl.AccessedStops[stop]; !exists {
		sl.AccessedStops[stop] = true
	}

	sl.Path = append(sl.Path, stop)
}

// HasVisited returns true if the stop has been visited
func (sl StopList) HasVisited(stop *graph.Stop) bool {
	_, visited := sl.AccessedStops[stop]
	return visited
}

// HasVisitedAllStops returns true if all stops have been visited
func (sl StopList) HasVisitedAllStops() bool {
	return len(sl.AccessedStops) == len(sl.allStops)
}
