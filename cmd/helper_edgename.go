package cmd

import "strings"

func makeStopNameValid(stopName string) string {
	finalStr := ""
	for _, s := range strings.Split(stopName, " ") {
		finalStr += strings.ToUpper(s[0:1]) + strings.ToLower(s[1:]) + " "
	}

	return finalStr[:len(finalStr)-1] // remove space at end
}
