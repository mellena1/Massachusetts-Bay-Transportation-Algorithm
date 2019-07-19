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

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("A fatal error occurred: %s", err)
	}
	startTime := time.Date(2019, time.July, 18, 6, 0, 0, 0, loc)

	timeFunctions := getCubicSplineFuncs()

	// calculation.PlotCubicSplineFunc(timeFunctions["Riverside:Bowdoin"], "riverside-bowdoin.png")
	// calculation.PlotCubicSplineFunc(timeFunctions["Riverside:Braintree"], "riverside-braintree.png")
	// calculation.PlotAllCubicSplineFuncs(timeFunctions, "AllRoutes.png")

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

func getCubicSplineFuncs() calculation.CubicSplineFunctionsHolder {
	endpoints := calculation.GetEndpointStops()

	if len(endpoints) == 0 {
		log.Fatalf("No endpoints returned")
	}

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("A fatal error occurred: %s", err)
	}
	startTime := time.Date(2019, time.July, 18, 6, 0, 0, 0, loc)
	endTime := time.Date(2019, time.July, 19, 0, 0, 0, 0, loc)
	interval := time.Minute * 30

	cubicSplineCalc, err := calculation.NewCubicSpline(readAPIKey("apikey.secret"))
	if err != nil {
		log.Fatal(err)
	}

	edges, _ := calculation.ReadAPICalls("edgesAPICalls_Thursday.json")
	cubicSplines := cubicSplineCalc.MakeCubicSplineFunctionForAllEdges(endpoints, interval, startTime, endTime, edges)
	// calculation.WriteCubicSplineFunctionsToFile(cubicSplines, "cubicSplineFunctions.json")

	// cubicSplineCalc.SaveAPICalls(endpoints, interval, startTime, endTime, "edgesAPICalls.json")

	return cubicSplines
}
