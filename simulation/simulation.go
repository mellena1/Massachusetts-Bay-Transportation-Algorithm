package simulation

import (
	"sync"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/graph"
)

// Simulation holds all info about a running simulation
type Simulation struct {
	Channel     chan SimData
	Data        SimData
	CurrentlyAt *graph.Stop
	GoingTo     *graph.Stop
	WG          *sync.WaitGroup
	// Vehicle     *Vehicle
}

func (s *Simulation) ArriveAtStation(arrivedAt *graph.Stop) {
	s.CurrentlyAt = arrivedAt
	s.Data.Stops.ArriveAtStop(arrivedAt)
}

func (s *Simulation) spawnNewSimulations(curStop *graph.Stop) {
	newSims := []Simulation{}
	for _, nextStopID := range curStop.Edges {
		nextStop := graph.StopMap[nextStopID]
		if !s.doneWithPath(nextStop) {
			newSim := Simulation{
				Channel:     s.Channel,
				Data:        s.Data,
				CurrentlyAt: curStop,
				GoingTo:     nextStop,
				WG:          s.WG,
			}
			newSim.Data.Stops = CloneStopList(s.Data.Stops)
			newSims = append(newSims, newSim)
		}
	}

	for _, newSim := range newSims {
		s.WG.Add(1)
		go newSim.Run()
	}
}

func (s *Simulation) doneWithPath(stop *graph.Stop) bool {
	visited := make(map[*graph.Stop]bool)
	// visited[stop] = true

	var dfs func(stop *graph.Stop) bool
	dfs = func(stop *graph.Stop) bool {
		visited[stop] = true
		if !s.Data.Stops.HasVisited(stop) {
			return false
		}
		for _, nextStopID := range stop.Edges {
			nextStop := graph.StopMap[nextStopID]
			if visited[nextStop] {
				continue
			}
			if !dfs(nextStop) {
				return false
			}
		}
		return true
	}

	return dfs(stop)
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
			s.ArriveAtStation(curStation)

			if s.Data.Stops.HasVisitedAllStops() {
				s.Channel <- s.Data
				return
			}

			if curStation.IsIntersection() {
				s.spawnNewSimulations(curStation)
				return
			}
			s.GoingTo = graph.StopMap[curStation.Edges[0]]
		}
	}
}
