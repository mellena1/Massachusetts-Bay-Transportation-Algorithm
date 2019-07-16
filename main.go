package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/calculation"
)

func main() {
	endpoints := calculation.GetEndpointStops()

	if len(endpoints) == 0 {
		log.Fatalf("No endpoints returned")
	}

	fmt.Println(len(endpoints))

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("A fatal error occurred: %s", err)
	}
	startTime := time.Date(2019, time.July, 18, 6, 0, 0, 0, loc)

	// timeFunctions, err := calculation.ReadLagrangeFunctionsFromFile("lagrangeFunctions.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	timeFunctions := getLagrangeFuncs()

	// dur := calculation.GetDurationForEdgeFromLagrange(timeFunctions["Riverside:Bowdoin"], time.Date(2019, time.July, 18, 6, 1, 0, 0, loc))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(dur.Hours())

	// calculation.PlotLagrangeFunc(timeFunctions["Riverside:Bowdoin"], "riverside-bowdoin.png")
	// calculation.PlotLagrangeFunc(timeFunctions["Riverside:Braintree"], "riverside-braintree.png")

	calc, err := calculation.NewCalculator(timeFunctions)
	if err != nil {
		log.Fatal(err)
	}
	route, duration := calc.FindBestRoute(endpoints, startTime)

	fmt.Printf("Trip Duration: %v\n", duration)

	for i, stop := range route {
		fmt.Printf("%d: %s\n", i, stop.Name)
	}
}

func readAPIKey(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Can't read api key: %v", err)
	}
	return string(data)
}

func getLagrangeFuncs() calculation.LagrangeFunctionsHolder {
	endpoints := calculation.GetEndpointStops()

	if len(endpoints) == 0 {
		log.Fatalf("No endpoints returned")
	}

	fmt.Println(len(endpoints))

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("A fatal error occurred: %s", err)
	}
	startTime := time.Date(2019, time.July, 18, 6, 0, 0, 0, loc)
	endTime := time.Date(2019, time.July, 19, 0, 0, 0, 0, loc)
	interval := time.Minute * 30

	lagrangeCalc, err := calculation.NewLagrange(readAPIKey("apikey.secret"))
	if err != nil {
		log.Fatal(err)
	}

	edges, _ := calculation.ReadAPICalls("edgesAPICalls.json")
	lagranges := lagrangeCalc.MakeLagrangeFunctionForAllEdges(endpoints, interval, startTime, endTime, edges)
	// calculation.WriteLangrageFunctionsToFile(lagranges, "lagrangeFunctionsCubicSplines.json")

	// lagrangeCalc.SaveAPICalls(endpoints, interval, startTime, endTime, "edgesAPICalls.json")

	return lagranges
}
