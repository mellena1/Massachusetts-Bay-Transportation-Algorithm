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
	for _, nextStopID := range curStop.Edges {
		nextStop := graph.StopMap[nextStopID]
		if !s.doneWithPath(nextStop, curStop) {
			newSim := Simulation{
				Channel:     s.Channel,
				Data:        s.Data,
				CurrentlyAt: curStop,
				GoingTo:     nextStop,
				WG:          s.WG,
			}
			newSim.Data.Stops = CloneStopList(s.Data.Stops)
			go newSim.Run()
		}
	}
}

func (s *Simulation) doneWithPath(stop *graph.Stop, lastStop *graph.Stop) bool {
	if !s.Data.Stops.HasVisited(stop) {
		return false
	}
	for _, nextStop := range stop.Edges {
		if graph.StopMap[nextStop] == lastStop {
			continue
		}
		if !s.doneWithPath(graph.StopMap[nextStop], stop) {
			return false
		}
	}
	return true
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
	s.WG.Add(1)
	for {
		// TODO:
		// if s.Vehicle == nil {

		// }

		// TODO: if train now at dest station
		if true {
			curStation := s.GoingTo
			s.ArriveAtStation(curStation)

			if s.Data.Stops.HasVisitedAllStops() {
				break
			}

			if curStation.IsIntersection() {
				s.spawnNewSimulations(curStation)
				s.WG.Done()
				return
			}
			s.GoingTo = graph.StopMap[curStation.Edges[0]]
		}
	}

	s.Channel <- s.Data
	s.WG.Done()
}
