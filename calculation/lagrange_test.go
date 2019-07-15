package calculation

import (
	"testing"
	"time"
)

func Test_lagrangeUnitFromTime(t *testing.T) {
	tests := []struct {
		timeString        string
		expectedfloatTime float64
	}{
		{"02 Jan 06 15:04 MST", 904},
		{"02 Jul 10 12:07 MST", 727},
		{"11 Aug 19 03:04 MST", 184},
	}

	for _, test := range tests {
		timeObj, _ := time.Parse(time.RFC822, test.timeString)
		actual := lagrangeUnitFromTime(timeObj)
		equals(t, test.expectedfloatTime, actual)
	}
}

func Test_lagrangeUnitFromDuration(t *testing.T) {
	tests := []struct {
		durString         string
		expectedfloatTime float64
	}{
		{"2h0m0s", 120},
		{"2h1m0s", 121},
		{"13h5m0s", 785},
	}

	for _, test := range tests {
		dur, _ := time.ParseDuration(test.durString)
		actual := lagrangeUnitFromDuration(dur)
		equals(t, test.expectedfloatTime, actual)
	}
}

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
