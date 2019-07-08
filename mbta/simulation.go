package mbta

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func main() error {
	err := importGraph("graph.json")
	if err != nil {
		return err
	}

	return nil
}

func importGraph(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Could not load file %s", filename)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Could not read bytes from %s", filename)
	}

	type tempStop struct {
		ID    string   `json:"id"`
		Name  string   `json:"name"`
		Edges []string `json:"edges"`
	}

	var tempStops []tempStop
	json.Unmarshal(bytes, &tempStops)

	return nil
}
