package simulation

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type routeData struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Path      []string  `json:"path"`
}

// ExportRoute stores route data in a text file
func ExportRoute(data SimData, filename string) error {
	path := make([]string, len(data.Stops.Path))
	for i, stop := range data.Stops.Path {
		path[i] = stop.ID
	}

	route := routeData{
		StartTime: data.StartTime,
		EndTime:   data.EndTime,
		Path:      path,
	}

	jsonData, err := json.MarshalIndent(route, "", "\t")
	if err != nil {
		return err
	}

	dir := "routes/" + route.StartTime.Format("Mon-01:01:2000") + "/"

	os.MkdirAll(dir, 0644)
	ioutil.WriteFile(dir+route.StartTime.Format("1:15PM")+": "+filename, jsonData, 0644)

	return nil
}
