package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/calculation"
	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"
	"github.com/spf13/cobra"
)

var reconstructDataFile string
var reconstructAPIKeyFile string

func init() {
	reconstructCmd.Flags().StringVarP(&reconstructDataFile, "input", "i", "", "File to get data from (required)")
	reconstructCmd.MarkFlagRequired("input")
	reconstructCmd.Flags().StringVarP(&reconstructAPIKeyFile, "apikey", "a", "apikey.secret", "The file containing the google maps api key")

	rootCmd.AddCommand(reconstructCmd)
}

var reconstructCmd = &cobra.Command{
	Use:   "reconstruct",
	Short: "Take calculated data and reconstruct route",
	Long:  "Take calculated data and reconstruct route",
	Run:   reconstructFunc,
}

func reconstructFunc(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(string(datacollection.SpecialEdgesFileWithLocationData)); err != nil {
		datacollection.GetStopCoordinatesForGoogleAPI()
	}
	if _, err := os.Stat(string(datacollection.StopLocations)); err != nil {
		datacollection.GetStopCoordinatesForGoogleAPI()
	}

	results, err := readResults(reconstructDataFile)
	if err != nil {
		fmt.Printf("Error reading results: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Pick a route to reconstruct:")
	i := 1
	keys := make([]time.Time, len(results))
	for t := range results {
		keys[i-1] = t
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})
	for i, t := range keys {
		result := results[t]
		fmt.Printf("%d. %s - %v - %s\n", i+1, t.Format(time.Kitchen), result.Duration, calculation.PrintStops(result.Route))
	}

	reader := bufio.NewReader(os.Stdin)
	var selectedNum int
	for {
		fmt.Print("-> ")
		input, _ := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)

		var err error
		selectedNum, err = strconv.Atoi(input)
		if err != nil || selectedNum < 1 || selectedNum > len(keys) {
			fmt.Println("please enter a valid number.")
			continue
		}
		break
	}

	startTime := keys[selectedNum-1]
	result := results[startTime]

	datacollection.ReconstructRoute(result.Route, startTime, reconstructAPIKeyFile)
}
