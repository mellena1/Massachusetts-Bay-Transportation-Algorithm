package graph

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

// ReadSpecialEdgesFromFile reads in the special edges from a file (probably special_edges.json)
func ReadSpecialEdgesFromFile(filename string) (SpecialEdges, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var specialEdges SpecialEdges
	err = json.Unmarshal(data, &specialEdges)

	return specialEdges, err
}
