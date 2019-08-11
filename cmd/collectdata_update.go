package cmd

import (
	"fmt"
	"os"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"
	"github.com/spf13/cobra"
)

var collectDataUpdateInputFile string
var collectDataUpdateStopAName string
var collectDataUpdateStopBName string

func init() {
	collectDataUpdateCmd.Flags().StringVarP(&collectDataUpdateInputFile, "input", "i", "", "File to update (required)")
	collectDataUpdateCmd.MarkFlagRequired("input")
	collectDataUpdateCmd.Flags().StringVar(&collectDataUpdateStopAName, "stopA", "", "Beginning stop of edge (required)")
	collectDataUpdateCmd.MarkFlagRequired("stopA")
	collectDataUpdateCmd.Flags().StringVar(&collectDataUpdateStopBName, "stopB", "", "Ending stop of edge (required)")
	collectDataUpdateCmd.MarkFlagRequired("stopB")

	collectDataCmd.AddCommand(collectDataUpdateCmd)
}

var collectDataUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Given an existing file, re-run data for a particular route",
	Long:  "Given an existing file, re-run data for a particular route",
	Run:   collectDataUpdateFunc,
}

func collectDataUpdateFunc(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(string(datacollection.SpecialEdgesFileWithLocationData)); err != nil {
		datacollection.GetStopCoordinatesForGoogleAPI()
	}
	if _, err := os.Stat(string(datacollection.StopLocations)); err != nil {
		datacollection.GetStopCoordinatesForGoogleAPI()
	}

	var stopA, stopB *datacollection.Stop

	stops, err := datacollection.ImportStopsFromFile(datacollection.StopLocations)
	if err != nil {
		fmt.Printf("Failed to import stop location data: %s", err)
		os.Exit(1)
	}

	collectDataUpdateStopAName = makeStopNameValid(collectDataUpdateStopAName)
	collectDataUpdateStopBName = makeStopNameValid(collectDataUpdateStopBName)

	for _, stop := range stops {
		if stop.Name == collectDataUpdateStopAName {
			stopA = stop
		}
		if stop.Name == collectDataUpdateStopBName {
			stopB = stop
		}
	}
	if stopA == nil {
		fmt.Printf("Invalid stop name for stopA")
		os.Exit(1)
	}
	if stopB == nil {
		fmt.Printf("Invalid stop name for stopB")
		os.Exit(1)
	}

	edgeData, err := datacollection.ImportEdgeData(collectDataUpdateInputFile)
	if err != nil {
		fmt.Printf("failed reading in edge data: %s", err)
		os.Exit(1)
	}

	parsedDate, err := datacollection.GetDateFromEdgeDataFilename(collectDataUpdateInputFile)
	if err != nil {
		fmt.Printf("invalid date in edgedata file. Please use format yyyy-mm-dd. %s", err)
		os.Exit(1)
	}

	startTime := parsedDate.Add(collectDataStartTimeOfDay)
	endTime := parsedDate.Add(collectDataEndTimeOfDay)

	newEdges := datacollection.GetTransitDataForAnEdgeWithGoogleAPI(stopA, stopB, startTime, endTime, collectDataInterval, collectDataAPIKeyFile)

	for k, v := range newEdges {
		edgeData[k] = v
	}

	err = datacollection.WriteEdgeDataToFile(edgeData, startTime, collectDataInterval)
	if err != nil {
		fmt.Printf("Failed to write to edgedata file. %s", err)
		os.Exit(1)
	}
}
