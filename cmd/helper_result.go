package cmd

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"
)

// Result stores calculated data from a calculation run
type Result struct {
	Route    []datacollection.Stop
	Duration time.Duration
}

// Results stores all results based on their start times
type Results map[time.Time]Result

func readResults(filename string) (Results, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	results := make(Results)
	err = json.Unmarshal(data, &results)
	return results, err
}
