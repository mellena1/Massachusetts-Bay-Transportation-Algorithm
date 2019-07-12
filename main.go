package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

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

	graph.InitPackage("graph copy.json")

	dataChannel := make(chan simulation.SimData)
	var wg sync.WaitGroup
	endpoints := getEndpoints()

	for _, stop := range endpoints {
		sim := simulation.Simulation{
			Channel: dataChannel,
			Data: simulation.SimData{
				StartTime: time.Now(),
				Stops:     simulation.NewStopList(graph.StopList, stop),
			},
			CurrentlyAt: stop,
			GoingTo:     graph.StopMap[stop.Edges[0]],
			WG:          &wg,
		}
		wg.Add(1)
		go sim.Run()
	}
	go func() {
		wg.Wait()
		close(dataChannel)
	}()

	x := 0
	for data := range dataChannel {
		x++
		err := simulation.ExportRoute(data, "test/route"+strconv.Itoa(x))
		if err != nil {
			log.Fatalf("Could not export route.")
		}
	}
	fmt.Println(x)
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
