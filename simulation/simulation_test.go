package simulation

import (
	"testing"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/graph"
)

func Test_arriveAtStop(t *testing.T) {
	graph.InitPackage("testdata/graph_nocycle.json")

	sim := Simulation{
		LastStop:    graph.StopList[0],
		CurrentlyAt: graph.StopList[4],
		GoingTo:     graph.StopList[1],
		Data: SimData{
			Stops: &StopList{AccessedStops: map[*graph.Stop]bool{graph.StopList[0]: true, graph.StopList[4]: true}, allStops: graph.StopList},
		},
	}

	sim.ArriveAtStop(sim.GoingTo)
	equals(t, map[*graph.Stop]bool{graph.StopList[0]: true, graph.StopList[4]: true, graph.StopList[1]: true}, sim.Data.Stops.AccessedStops)
	equals(t, graph.StopList[4], sim.LastStop)
	equals(t, graph.StopList[1], sim.CurrentlyAt)
}

func Test_doneWithPath(t *testing.T) {
	tests := []struct {
		graphFile          string
		nextStopIndex      int
		visitedStopIndices []int
		expected           bool
	}{
		{"testdata/graph_cycle.json", 1, []int{}, false},
		{"testdata/graph_cycle.json", 1, []int{1}, false},
		{"testdata/graph_cycle.json", 4, []int{0, 1, 2, 3, 4}, true},
		{"testdata/graph_nocycle.json", 1, []int{}, false},
		{"testdata/graph_nocycle.json", 1, []int{1}, false},
		{"testdata/graph_nocycle.json", 4, []int{0, 1, 2, 3, 4}, true},
	}

	for _, test := range tests {
		graph.InitPackage(test.graphFile)

		sim := Simulation{
			Data: SimData{
				Stops: &StopList{AccessedStops: make(map[*graph.Stop]bool), allStops: graph.StopList},
			},
		}
		for _, index := range test.visitedStopIndices {
			sim.Data.Stops.ArriveAtStop(graph.StopList[index])
		}

		//actual := sim.doneWithPath(graph.StopList[test.nextStopIndex])
		//equals(t, test.expected, actual)
	}
}
