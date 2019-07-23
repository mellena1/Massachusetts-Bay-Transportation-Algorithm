package cmd

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"
	"github.com/spf13/cobra"
)

var collectDataViewInputFile string
var collectDataViewStopAName string
var collectDataViewStopBName string

func init() {
	collectDataViewCmd.Flags().StringVarP(&collectDataViewInputFile, "input", "i", "", "File to update (required)")
	collectDataViewCmd.MarkFlagRequired("input")
	collectDataViewCmd.Flags().StringVar(&collectDataViewStopAName, "stopA", "", "Beginning stop of edge (required)")
	collectDataViewCmd.Flags().StringVar(&collectDataViewStopBName, "stopB", "", "Ending stop of edge (required)")

	collectDataCmd.AddCommand(collectDataViewCmd)
}

var collectDataViewCmd = &cobra.Command{
	Use:   "view",
	Short: "Given an existing file, view the timings",
	Long:  "Given an existing file, view the timings",
	Run:   collectDataViewFunc,
}

func collectDataViewFunc(cmd *cobra.Command, args []string) {
	edgeData, err := datacollection.ImportEdgeData(collectDataViewInputFile)
	if err != nil {
		fmt.Printf("failed reading in edge data: %s", err)
		os.Exit(1)
	}

	if (collectDataViewStopAName != "" && collectDataViewStopBName == "") || (collectDataViewStopAName == "" && collectDataViewStopBName != "") {
		fmt.Println("Must specify both stop names if specifying one")
		os.Exit(1)
	}

	if collectDataViewStopAName != "" && collectDataViewStopBName != "" {
		if timings, ok := edgeData[datacollection.GetEdgeKey(collectDataViewStopAName, collectDataViewStopBName)]; ok {
			printEdgeTimings(timings)
		} else {
			fmt.Println("Stop names not found in data")
			os.Exit(1)
		}
		if timings, ok := edgeData[datacollection.GetEdgeKeyWalking(collectDataViewStopAName, collectDataViewStopBName)]; ok {
			fmt.Printf("--- Walking ---")
			printEdgeTimings(timings)
		}
	} else {
		for edge, timings := range edgeData {
			fmt.Printf("--- %s\n", edge)
			printEdgeTimings(timings)
			fmt.Println()
		}
	}
}

func printEdgeTimings(timings datacollection.EdgeTimes) {
	keys := make([]time.Time, len(timings))
	i := 0
	for unixTime := range timings {
		startTime := time.Unix(unixTime, 0)
		keys[i] = startTime
		i++
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	for _, startTime := range keys {
		fmt.Printf("Time: %s - Duration: %s\n", startTime.Format(time.Kitchen), timings[startTime.Unix()])
	}
}
