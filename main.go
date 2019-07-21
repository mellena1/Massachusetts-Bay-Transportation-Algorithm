package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/calculation"
	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"
)

func main() {
	// datacollection.GetStopCoordinatesForGoogleAPI()

	endpoints, err := datacollection.ImportStopsFromFileNonePointer(datacollection.Endpoints)
	if err != nil {
		log.Fatalf("can't read endpoints: %s", err)
	}

	// if len(endpoints) == 0 {
	// 	log.Fatalf("No endpoints returned")
	// }

	// loc, _ := time.LoadLocation("America/New_York")
	// startTime := time.Date(2019, time.July, 24, 6, 0, 0, 0, loc)
	// endTime := time.Date(2019, time.July, 25, 0, 0, 0, 0, loc)
	// interval := time.Minute * 30

	// datacollection.GetTransitDataWithGoogleAPI(startTime, endTime, interval)

	// calculation.PlotCubicSplineFunc(timeFunctions["Riverside:Bowdoin"], "riverside-bowdoin.png")
	// calculation.PlotCubicSplineFunc(timeFunctions["Riverside:Braintree"], "riverside-braintree.png")
	// calculation.PlotAllCubicSplineFuncs(timeFunctions, "AllRoutes.png")

	edgeData, err := datacollection.ImportEdgeData("datacollection/EdgeData StartTime:1563962400 Interval:30m0s.json")
	if err != nil {
		log.Fatalf("failed reading in edge data: %s", err)
	}

	loc, _ := time.LoadLocation("America/New_York")
	startTime := time.Date(2019, time.July, 24, 6, 0, 0, 0, loc)
	lastTime := time.Date(2019, time.July, 24, 19, 0, 0, 0, loc)
	interval := time.Hour

	type result struct {
		route    []calculation.Stop
		duration time.Duration
	}
	results := make(map[time.Time]result)
	routesLock := sync.Mutex{}

	numberOfRunners := struct {
		num int
		mu  sync.Mutex
	}{}

	numberOfCores := runtime.NumCPU()

	wg := sync.WaitGroup{}
	for curTime := startTime; curTime.Before(lastTime) || curTime.Equal(lastTime); curTime = curTime.Add(interval) {
		numberOfRunners.mu.Lock()
		numberOfRunners.num++
		numberOfRunners.mu.Unlock()

		go func(t time.Time) {
			wg.Add(1)

			calc, err := calculation.NewCalculator(edgeData, time.Date(2019, time.July, 25, 0, 0, 0, 0, loc))
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Starting: %v", t)
			route, duration := calc.FindBestRoute(endpoints, t)
			routesLock.Lock()
			defer routesLock.Unlock()
			results[t] = result{route: route, duration: duration}
			log.Printf("Done with: %v", t)

			wg.Done()
			numberOfRunners.mu.Lock()
			numberOfRunners.num--
			numberOfRunners.mu.Unlock()
		}(curTime)

		for {
			numberOfRunners.mu.Lock()
			num := numberOfRunners.num
			numberOfRunners.mu.Unlock()
			if num == numberOfCores {
				time.Sleep(time.Second)
			} else {
				break
			}
		}
	}
	wg.Wait()

	for t, res := range results {
		log.Printf("%v --- Duration: %v Route: %s", t, res.duration, calculation.PrintStops(res.route))
	}
	data, err := json.Marshal(&results)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("results.json", data, 0644)
}

func readAPIKey(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Can't read api key: %v", err)
	}
	return string(data)
}

// func omain() {
// 	endpoints := calculation.GetEndpointStops()

// 	if len(endpoints) == 0 {
// 		log.Fatalf("No endpoints returned")
// 	}

// 	loc, err := time.LoadLocation("America/New_York")
// 	if err != nil {
// 		log.Fatalf("A fatal error occurred: %s", err)
// 	}
// 	startTime := time.Date(2019, time.July, 18, 6, 0, 0, 0, loc)

// 	// timeFunctions, err := calculation.ReadCubicSplineFunctionsFromFile("cubicSplineFunctions.json")
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	timeFunctions := getCubicSplineFuncs()

// 	// dur := calculation.GetDurationForEdgeFromCubicSpline(timeFunctions["Riverside:Bowdoin"], time.Date(2019, time.July, 18, 6, 1, 0, 0, loc))
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// fmt.Println(dur.Hours())

// 	calculation.PlotCubicSplineFunc(timeFunctions["Riverside:Bowdoin"], "riverside-bowdoin.png")
// 	calculation.PlotCubicSplineFunc(timeFunctions["Riverside:Braintree"], "riverside-braintree.png")
// 	calculation.PlotAllCubicSplineFuncs(timeFunctions, "AllRoutes.png")

// 	calc, err := calculation.NewCalculator(timeFunctions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	route, duration := calc.FindBestRoute(endpoints, startTime)

// 	fmt.Printf("Trip Duration: %v\n", duration)

// 	for i, stop := range route {
// 		fmt.Printf("%d: %s\n", i, stop.Name)
// 	}
// }

// func getCubicSplineFuncs() calculation.CubicSplineFunctionsHolder {
// 	endpoints := calculation.GetEndpointStops()

// 	if len(endpoints) == 0 {
// 		log.Fatalf("No endpoints returned")
// 	}

// 	loc, err := time.LoadLocation("America/New_York")
// 	if err != nil {
// 		log.Fatalf("A fatal error occurred: %s", err)
// 	}
// 	startTime := time.Date(2019, time.July, 18, 6, 0, 0, 0, loc)
// 	endTime := time.Date(2019, time.July, 19, 0, 0, 0, 0, loc)
// 	interval := time.Minute * 30

// 	cubicSplineCalc, err := calculation.NewCubicSpline(readAPIKey("apikey.secret"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	edges, _ := calculation.ReadAPICalls("edgesAPICalls_Thursday.json")
// 	cubicSplines := cubicSplineCalc.MakeCubicSplineFunctionForAllEdges(endpoints, interval, startTime, endTime, edges)
// 	// calculation.WriteCubicSplineFunctionsToFile(cubicSplines, "cubicSplineFunctions.json")

// 	// cubicSplineCalc.SaveAPICalls(endpoints, interval, startTime, endTime, "edgesAPICalls.json")

// 	return cubicSplines
// }
