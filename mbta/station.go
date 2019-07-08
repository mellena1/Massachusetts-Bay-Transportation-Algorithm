package mbta

import (
	"fmt"

	mbtaClient "github.com/mellena1/mbta-v3-go/mbta"
)

// IStation defines an mbta train station
type IStation interface {
	IGetID
	ISetID
	IGetNextTrain
}

// Station represents an mbta train station
type Station struct {
	Node
	mbtaClient.Stop
}

// IGetID defines a function for returning a station's id
type IGetID interface {
	GetId() string
}

// GetID returns a station's id
func (s *Station) GetID() string {
	return s.Stop.ID
}

// ISetID defines a function to set a station's id
type ISetID interface {
	SetId(string)
}

// SetID updates a station's id
func (s *Station) SetID(id string) {
	s.Stop.ID = id
}

// IGetNextTrain defines a function that gets the next train to enter this station
type IGetNextTrain interface {
	GetNextTrain()
}

// GetNextTrain returns the next train to enter this station
func (s *Station) GetNextTrain() error {
	client := mbtaClient.NewClient(mbtaClient.ClientConfig{})

	vehicles, _, err := client.Vehicles.GetAllVehicles(&mbtaClient.GetAllVehiclesRequestConfig{})
	if err != nil {
		return err
	}

	fmt.Println(vehicles)

	return nil
}
