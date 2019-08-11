package calculation

import (
	"testing"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"
)

func Test_removeIndex(t *testing.T) {
	list := []datacollection.Stop{
		datacollection.Stop{
			Name: "Andrew",
		},
		datacollection.Stop{
			Name: "Brad",
		},
		datacollection.Stop{
			Name: "Charles",
		},
		datacollection.Stop{
			Name: "Sam",
		},
	}

	expectedList := []datacollection.Stop{
		datacollection.Stop{
			Name: "Andrew",
		},
		datacollection.Stop{
			Name: "Brad",
		},
		datacollection.Stop{
			Name: "Charles",
		},
		datacollection.Stop{
			Name: "Sam",
		},
	}

	expectedNewList := []datacollection.Stop{
		datacollection.Stop{
			Name: "Andrew",
		},
		datacollection.Stop{
			Name: "Brad",
		},
		datacollection.Stop{
			Name: "Sam",
		},
	}

	newList := removeIndex(2, list)

	equals(t, expectedNewList, newList)
	equals(t, expectedList, list)
}
