package datacollection

import "strconv"

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
