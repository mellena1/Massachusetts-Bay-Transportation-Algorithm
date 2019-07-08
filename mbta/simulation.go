package mbta

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() error {
	stops, err := importGraph()
	if err != nil {
		return err
	}

	fmt.Println(stops)

	return nil
}

func importGraph() (*[]Stop, error) {
	jsonFile, err := os.Open("graph.json")
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var stops []Stop
	json.Unmarshal(data, &stops)

	return &stops, nil
}
