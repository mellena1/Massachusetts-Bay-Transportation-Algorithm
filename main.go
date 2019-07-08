package main

import (
	"fmt"
	"sync"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/graph"
	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/simulation"
)

func main() {
	/*
		main thread:
			make data channel for sims to send done data back
			get the starting stops from the graph
			for range starting stops:
				check when trains come in
				if train:
					make simulation struct
					go sim.Run()

			for range data channel:
				if something:
					write data to db/file
	*/

	dataChannel := make(chan simulation.SimData)
	var wg sync.WaitGroup
	endpoints := getEndpoints()

	for _, stop := range endpoints {
		sim := simulation.Simulation{
			Channel: dataChannel,
			Data: simulation.SimData{
				Stops: simulation.NewStopList(graph.StopList),
			},
			CurrentlyAt: stop,
			GoingTo:     graph.StopMap[stop.Edges[0]],
			WG:          &wg,
		}
		go sim.Run()
	}
	// wg.Wait()

	// x := 0
	// for {
	// 	_, ok := <-dataChannel
	// 	if !ok {
	// 		break
	// 	}
	// 	x++
	// }
	// fmt.Println(x)

	n := 0
	for range dataChannel {
		n++
		fmt.Println(n)
		panic("hi")
	}
}

func getEndpoints() []*graph.Stop {
	endpoints := []*graph.Stop{}
	for _, s := range graph.StopList {
		if s.IsEndpoint() {
			endpoints = append(endpoints, s)
		}
	}
	return endpoints
}
