package mbta

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Stop represents an mbta stop
type Stop struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Edges []*Stop `json:"edges"`
}

// IsEndpoint returns true if the stop is on the end of a line
func (s *Stop) IsEndpoint() bool {
	return len(s.Edges) == 1
}

// IsIntersection returns true if stop is an intersecting stop
func (s *Stop) IsIntersection() bool {
	return len(s.Edges) > 1
}

// UnmarshalJSON customer unmarshaller for the graph
func (s *Stop) UnmarshalJSON(b []byte) error {
	file, err := os.Open("graph.json")
	if err != nil {
		log.Fatal("Could not load file graph.json")
	}

	type tempStop struct {
		ID    string    `json:"id"`
		Name  string    `json:"name"`
		Edges []*string `json:"edges"`
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Could not read bytes from file graph.json")
	}

	var tempStops []tempStop
	json.Unmarshal(data, &stops)

	return nil
}
