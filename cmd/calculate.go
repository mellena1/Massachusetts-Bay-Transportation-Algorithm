package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/calculation"
	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"
	"github.com/spf13/cobra"
)

var calculateInputFile string
var calculateOutputFile string

func init() {
	calculateCmd.Flags().StringVarP(&calculateInputFile, "input", "i", "", "File read the edge data from")
	calculateCmd.MarkFlagRequired("input")

	calculateCmd.Flags().StringVarP(&calculateOutputFile, "output", "o", "results.json", "File to save the results to")

	rootCmd.AddCommand(calculateCmd)
}

var calculateCmd = &cobra.Command{
	Use:   "calculate",
	Short: "Given a data file, find the fastest time for that day",
	Long:  "Given a data file, find the fastest time for that day",
	Run:   calculateFunc,
}

func calculateFunc(cmd *cobra.Command, args []string) {
	endpoints, err := datacollection.ImportStopsFromFileNonePointer(datacollection.Endpoints)
	if err != nil {
		log.Fatalf("can't read endpoints: %s", err)
	}

	edgeData, err := datacollection.ImportEdgeData(calculateInputFile)
	if err != nil {
		log.Fatalf("failed reading in edge data: %s", err)
	}

	firstRouteStartTime, _ := time.Parse(time.Kitchen, "6:00AM")
	lastRouteStartTime, _ := time.Parse(time.Kitchen, "7:00PM")
	interval := time.Hour

	// Cubic spline data ends at midnight, don't try calculating after that
	latestRouteTime, _ := time.Parse(time.Kitchen, "12:00AM")
	latestRouteTime = latestRouteTime.Add(time.Hour * 24)

	type Result struct {
		Route    []calculation.Stop
		Duration time.Duration
	}
	results := make(map[time.Time]Result)
	routesLock := sync.Mutex{}

	numberOfRunners := struct {
		num int
		mu  sync.Mutex
	}{}
	numberOfCores := runtime.NumCPU()

	wg := sync.WaitGroup{}
	for curTime := firstRouteStartTime; curTime.Before(lastRouteStartTime) || curTime.Equal(lastRouteStartTime); curTime = curTime.Add(interval) {
		numberOfRunners.mu.Lock()
		numberOfRunners.num++
		numberOfRunners.mu.Unlock()

		go func(t time.Time) {
			wg.Add(1)

			calc, err := calculation.NewCalculator(edgeData, latestRouteTime)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Starting: %s", t.Format(time.Kitchen))
			route, duration := calc.FindBestRoute(endpoints, t)

			routesLock.Lock()
			results[t] = Result{Route: route, Duration: duration}
			routesLock.Unlock()

			log.Printf("Done with: %s", t.Format(time.Kitchen))
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
		fmt.Printf("%v --- Duration: %v Route: %s\n", t, res.Duration, calculation.PrintStops(res.Route))
	}
	data, err := json.Marshal(&results)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(calculateOutputFile, data, 0644)
}
