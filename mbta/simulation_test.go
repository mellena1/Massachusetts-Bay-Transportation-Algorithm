package mbta

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test_importGraph(t *testing.T) {
	jsonFile, err := os.Open("graph.json")
	if err != nil {
		t.FailNow()
	}

	expected, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		t.FailNow()
	}

	stops, err := importGraph()
	if err != nil {
		t.FailNow()
	}

	actual, err := json.MarshalIndent(stops, "", "    ")
	if err != nil {
		t.FailNow()
	}

	fmt.Printf("act: %v\n\nexp: %v\n", string(actual), string(expected))

	equals(t, actual, expected)
}
