package mbtahelper

import "github.com/mellena1/mbta-v3-go/mbta"

// MBTAHelper gives helper methods for doing stuff with the MBTA API
type MBTAHelper struct {
	Client *mbta.Client
}

// NewMBTAHelper returns a new client given an API key
func NewMBTAHelper(APIKey string) *MBTAHelper {
	return &MBTAHelper{
		Client: mbta.NewClient(mbta.ClientConfig{
			APIKey: APIKey,
		}),
	}
}
