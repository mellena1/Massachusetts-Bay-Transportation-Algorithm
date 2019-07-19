package datacollection

import (
	"encoding/json"
	"io/ioutil"
)

// SpecialEdge is a struct holding the MiddleEdgeID
type SpecialEdge struct {
	MiddleEdgeID string `json:"middle-stop"`
}

// SpecialEdges holds which edges need a stop in between them (and can also walk between them)
type SpecialEdges map[string]SpecialEdge

// SpecialEdgeFiles enum for files containing special edges
type SpecialEdgeFiles string

const (
	// SpecialEdgesFile the file that stores the edges with multiple ways to go
	SpecialEdgesFile SpecialEdgeFiles = "datacollection/special_edges.json"
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
