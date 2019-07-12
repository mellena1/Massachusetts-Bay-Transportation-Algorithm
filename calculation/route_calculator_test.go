package calculation

import (
	"fmt"
	"testing"
	"time"
)

func Test_FindBestRoute(t *testing.T) {
	stops := []string{"Haymarket", "North Station", "Oak Grove", "Lechmere"}
	FindBestRoute(stops, time.Now())
	fmt.Println(numberOfRoutes)
	t.FailNow()
}
