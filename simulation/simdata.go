package simulation

import "time"

// SimData holds info about a simulation, like start and end time and
// which stops have been visited
type SimData struct {
	StartTime time.Time
	EndTime   time.Time
	Stops     *StopList
}
