package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/calculation"
)

func main() {
	endpoints := calculation.GetEndpointStops()

	if len(endpoints) == 0 {
		log.Fatalf("No endpoints returned")
	}

	fmt.Println(endpoints)

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("A fatal error occurred: %s", err)
	}
	startTime := time.Date(2019, time.July, 18, 9, 0, 0, 0, loc)

	route, duration := calculation.FindBestRoute(endpoints, startTime)

	fmt.Printf("Trip Duration: %v\n", duration)

	for i, stop := range route {
		fmt.Printf("%d: %s\n", i, stop.Name)
	}
}
