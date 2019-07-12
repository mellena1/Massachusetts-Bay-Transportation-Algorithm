package calculation

import (
	"fmt"
	"testing"
)

func Test_FindBestRoute(t *testing.T) {
	stops := []string{"stop1", "stop2", "stop3", "stop4", "stop5"}
	FindBestRoute(stops)
	fmt.Println(numberOfRoutes)
	t.FailNow()
}
