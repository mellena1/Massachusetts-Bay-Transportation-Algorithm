package datacollection

import (
	"encoding/json"
	"io/ioutil"
)

// SpecialEdges holds which edges need a stop in between them (and can also walk between them). The key is the edge and the value is the middle stop.
type SpecialEdges map[string]*Stop

// SpecialEdgeFiles enum for files containing special edges
type SpecialEdgeFiles string

const (
	// SpecialEdgesFile the file that stores the edges with multiple ways to go
	SpecialEdgesFile SpecialEdgeFiles = "datacollection/special_edges.json"
	// SpecialEdgesFileWithLocationData the file that stores the edges with multiple ways to go and the middle stop locations
	SpecialEdgesFileWithLocationData SpecialEdgeFiles = "datacollection/special_edges_with_location_data.json"
)

// ReadSpecialEdgesFromFile reads in the special edges from a file (probably special_edges.json)
func ReadSpecialEdgesFromFile(filename SpecialEdgeFiles) (SpecialEdges, error) {
	data, err := ioutil.ReadFile(string(filename))
	if err != nil {
		return nil, err
	}
	var specialEdges SpecialEdges
	err = json.Unmarshal(data, &specialEdges)

	return specialEdges, err
}
