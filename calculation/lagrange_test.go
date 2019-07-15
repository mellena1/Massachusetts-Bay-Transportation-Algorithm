package calculation

import (
	"testing"
)

func Test_durationFromLagrangeUnit(t *testing.T) {
	tests := []struct {
		floatTime         float64
		expectedDurString string
	}{
		{120, "2h0m0s"},
		{121, "2h1m0s"},
		{785, "13h5m0s"},
	}

	for _, test := range tests {
		actual := durationFromLagrangeUnit(test.floatTime)
		equals(t, test.expectedDurString, actual.String())
	}
}
