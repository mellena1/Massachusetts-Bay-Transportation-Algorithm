package graph

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Stop represents an mbta stop
type Stop struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Edges []*string `json:"edges"`
}

// StopList list of stops
var StopList []Stop

// StopMap map of stop IDs to Stops
var StopMap map[string]Stop

func init() {
	var err error

	StopList, err = importGraph("graph.json")
	if err != nil {
		log.Fatalf("Could not load file graph.json")
	}

}

// IsEndpoint returns true if the stop is on the end of a line
func (s *Stop) IsEndpoint() bool {
	return len(s.Edges) == 1
}

// IsIntersection returns true if stop is an intersecting stop
func (s *Stop) IsIntersection() bool {
	return len(s.Edges) > 1
}

func importGraph(filename string) ([]Stop, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var stops []Stop
	json.Unmarshal(bytes, &stops)
	return stops, nil
}
