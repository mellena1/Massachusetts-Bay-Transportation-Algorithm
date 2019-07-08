package mbta

// Stop represents an mbta stop
type Stop struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	EdgeIDs []string `json:"edgeIDs"`
}

// IsEndpoint returns true if the stop is on the end of a line
func (s *Stop) IsEndpoint() bool {
	return len(s.EdgeIDs) == 1
}
