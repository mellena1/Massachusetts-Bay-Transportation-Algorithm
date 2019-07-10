package simulation

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/graph"
)

type routeData struct {
	StartTime time.Time     `json:"startTime"`
	EndTime   time.Time     `json:"endTime"`
	Path      []*graph.Stop `json:"path"`
}

// ExportRoute stores route data in a text file
func ExportRoute(data SimData, filename string) error {
	route := routeData{
		StartTime: data.StartTime,
		EndTime:   data.EndTime,
		Path:      data.Stops.Path,
	}

	jsonData, err := json.Marshal(route)
	if err != nil {
		return err
	}

	ioutil.WriteFile(filename, jsonData, 0644)

	return nil
}
