package graph

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Stop represents an mbta stop
type Stop struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Edges []string `json:"edges"`
}

// IsEndpoint returns true if the stop is on the end of a line
func (s *Stop) IsEndpoint() bool {
	return len(s.Edges) == 1
}

// IsIntersection returns true if stop is an intersecting stop
func (s *Stop) IsIntersection() bool {
	return len(s.Edges) > 1
}

// StopList list of stops
var StopList []*Stop

// StopMap map of stop IDs to Stops
var StopMap map[string]*Stop

// InitPackage must run this function before using the package
func InitPackage(filename string) {
	var err error

	StopList, StopMap, err = LoadGraphFile(filename)
	if err != nil {
		log.Fatalf("Could not load file \"%s\"", filename)
	}
}

// LoadGraphFile creates an array of all Stops given a json file along
// with a map with the ID of each stop as a string and the Stop as a
// pointer as the value
func LoadGraphFile(filename string) ([]*Stop, map[string]*Stop, error) {
	stopList, err := importGraph(filename)
	if err != nil {
		return nil, nil, err
	}

	stopMap := make(map[string]*Stop, len(StopList))

	for _, stop := range stopList {
		stopMap[stop.ID] = stop
	}

	return stopList, stopMap, nil
}

func importGraph(filename string) ([]*Stop, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var stops []*Stop
	json.Unmarshal(bytes, &stops)
	return stops, nil
}
