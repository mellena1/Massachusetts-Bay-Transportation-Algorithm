package mbta

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
