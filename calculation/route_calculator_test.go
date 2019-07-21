package calculation

import (
	"testing"
)

func Test_removeIndex(t *testing.T) {
	list := []Stop{
		Stop{
			Name: "Andrew",
		},
		Stop{
			Name: "Brad",
		},
		Stop{
			Name: "Charles",
		},
		Stop{
			Name: "Sam",
		},
	}

	expectedList := []Stop{
		Stop{
			Name: "Andrew",
		},
		Stop{
			Name: "Brad",
		},
		Stop{
			Name: "Charles",
		},
		Stop{
			Name: "Sam",
		},
	}

	expectedNewList := []Stop{
		Stop{
			Name: "Andrew",
		},
		Stop{
			Name: "Brad",
		},
		Stop{
			Name: "Sam",
		},
	}

	newList := removeIndex(2, list)

	equals(t, expectedNewList, newList)
	equals(t, expectedList, list)
}
