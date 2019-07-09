package simulation

import (
	"sync"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/graph"
)

// Simulation holds all info about a running simulation
type Simulation struct {
	Channel     chan SimData
	Data        SimData
	LastStop    *graph.Stop
	CurrentlyAt *graph.Stop
	GoingTo     *graph.Stop
	WG          *sync.WaitGroup
	// Vehicle     *Vehicle
}

// ArriveAtStop update simulation to being at the new Stop and add it to the path
func (s *Simulation) ArriveAtStop(arrivedAt *graph.Stop) {
	s.LastStop = s.CurrentlyAt
	s.CurrentlyAt = arrivedAt
	s.Data.Stops.ArriveAtStop(arrivedAt)
}

func (s *Simulation) createNewIntersectionSimulations() []*Simulation {
	newSims := []*Simulation{}
	curStop := s.CurrentlyAt

	for _, nextStopID := range curStop.Edges {
		nextStop := graph.StopMap[nextStopID]
		if nextStop == s.LastStop { // Don't go backwards to avoid cycle issues
			continue
		}
		if !s.doneWithPath(nextStop, curStop) {
			newSim := Simulation{
				Channel:     s.Channel,
				Data:        s.Data,
				LastStop:    s.LastStop,
				CurrentlyAt: curStop,
				GoingTo:     nextStop,
				WG:          s.WG,
			}
			newSim.Data.Stops = CloneStopList(s.Data.Stops)
			newSims = append(newSims, &newSim)
		}
	}

	return newSims
}

func (s *Simulation) doneWithPath(stop *graph.Stop, lastStop *graph.Stop) bool {
	visited := make(map[*graph.Stop]bool)

	var dfs func(stop *graph.Stop, lastStop *graph.Stop) ([]*graph.Stop, bool)
	dfs = func(stop *graph.Stop, lastStop *graph.Stop) ([]*graph.Stop, bool) {
		visited[stop] = true
		if !s.Data.Stops.HasVisited(stop) {
			pathTaken := []*graph.Stop{stop}
			return pathTaken, false
		}
		for _, nextStopID := range stop.Edges {
			nextStop := graph.StopMap[nextStopID]
			if nextStop == lastStop {
				continue
			}
			if visited[nextStop] {
				continue
			}
			if pathTaken, visited := dfs(nextStop, stop); !visited {
				return append(pathTaken, stop), false
			}
		}
		return nil, true
	}

	pathTaken, done := dfs(stop, lastStop)
	if done || isCycleInPath(pathTaken) {
		return true
	}

	return false
}

func isCycleInPath(path []*graph.Stop) bool {
	visited := make(map[*graph.Stop]bool)
	for _, stop := range path {
		if visited[stop] {
			return true
		}
		visited[stop] = true
	}
	return false
}

/*
	simulation thread:
			if hit all stations:
				send data to channel and end
			if train == nil:
				func: wait for a train to come at current station
			for:
				if train now at dest station:
					if intersection:
						make all possible simulations
						go sim.Run()
*/

// Run is the method to start the simulator
// make sure to run this in a goroutine like:
// go s.Run()
func (s *Simulation) Run() {
	defer s.WG.Done()
	for {
		// TODO:
		// if s.Vehicle == nil {

		// }

		// TODO: if train now at dest station
		if true {
			curStation := s.GoingTo
			s.ArriveAtStop(curStation)

			if s.Data.Stops.HasVisitedAllStops() {
				s.Channel <- s.Data
				return
			}

			if curStation.IsIntersection() {
				newSims := s.createNewIntersectionSimulations()
				for _, newSim := range newSims {
					s.WG.Add(1)
					go newSim.Run()
				}
				return
			}
			s.GoingTo = graph.StopMap[curStation.Edges[0]]
		}
	}
}
