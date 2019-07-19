package datacollection

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// StopFiles enum of all stop data files
type StopFiles string

const (
	// Endpoints the file that stores the endpoint names
	Endpoints StopFiles = "datacollection/endpoint_stops.json"
	// StopLocations the files that stores the locations of mbta stops
	StopLocations StopFiles = "datacollection/stop_locations.json"	
)

// ImportStopsFromFile imports a list of stop data
func ImportStopsFromFile(filename StopFiles) ([]*Stop, error) {
	file, err := os.Open(string(filename))
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
