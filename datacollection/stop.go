package datacollection

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)

// StopFiles enum of all stop data files
type StopFiles string

const (
	// Endpoints the file that stores the endpoint names
	Endpoints StopFiles = "datacollection/endpoint_stops.json"
	// StopLocations the files that stores the locations of mbta stops
	StopLocations StopFiles = "datacollection/stop_locations.json"
)

// Stop stores the coordinates of a stop and its name
type Stop struct {
	Name                   string `json:"name"`
	LongitudeCommaLatitude string `json:"longitude_comma_latitude,omitempty"`
}

// NewStop returns a new stop struct
func NewStop(name string, longitude string, latitude string) *Stop {
	return &Stop{
		Name:                   name,
		LongitudeCommaLatitude: latitude + "," + longitude,
	}
}

// SetLongitudeCommaLatitude sets the LongitudeCommaLatitude field given longitude and latitude floats
func (s *Stop) SetLongitudeCommaLatitude(longitude float64, latitude float64) {
	s.LongitudeCommaLatitude = strconv.FormatFloat(latitude, 'f', -1, 64) + "," + strconv.FormatFloat(longitude, 'f', -1, 64)
}

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

// ImportStopsFromFileNonePointer imports a list of stop data
func ImportStopsFromFileNonePointer(filename StopFiles) ([]Stop, error) {
	file, err := os.Open(string(filename))
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
