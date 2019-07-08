package mbta

import (
	"testing"
)

func Test_GetID(t *testing.T) {
	//t.FailNow()
}

func Test_SetID(t *testing.T) {
	//t.FailNow()
}

func Test_GetNextTrain(t *testing.T) {
	station := Station{}
	err := station.GetNextTrain()
	if err != nil {
		t.Fail()
	}
	t.FailNow()
}
