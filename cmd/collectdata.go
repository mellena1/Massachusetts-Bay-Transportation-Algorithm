package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"
	"github.com/spf13/cobra"
)

var collectDataDateString string
var collectDataAPIKeyFile string

const (
	collectDataStartTimeOfDay time.Duration = time.Hour * 6  // 6AM
	collectDataEndTimeOfDay   time.Duration = time.Hour * 24 // midnight next day
	collectDataInterval       time.Duration = time.Minute * 30
)

func init() {
	collectDataCmd.Flags().StringVarP(&collectDataDateString, "date", "d", "", "Date to collect data for. Format: yyyy-mm-dd (required)")
	collectDataCmd.MarkFlagRequired("date")
	collectDataCmd.PersistentFlags().StringVarP(&collectDataAPIKeyFile, "apikey", "a", "apikey.secret", "The file containing the google maps api key")

	rootCmd.AddCommand(collectDataCmd)
}

var collectDataCmd = &cobra.Command{
	Use:   "collectdata",
	Short: "Given a day, go collect data from google maps",
	Long:  "Given a day, go collect data from google maps",
	Run:   collectDataFunc,
}

func collectDataFunc(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(string(datacollection.SpecialEdgesFileWithLocationData)); err != nil {
		datacollection.GetStopCoordinatesForGoogleAPI()
	}
	if _, err := os.Stat(string(datacollection.StopLocations)); err != nil {
		datacollection.GetStopCoordinatesForGoogleAPI()
	}

	date, err := datacollection.ParseDateToEST(collectDataDateString)
	if err != nil {
		fmt.Printf("invalid date: %s; error: %s", collectDataDateString, err)
		os.Exit(1)
	}

	startTime := date.Add(collectDataStartTimeOfDay)
	endTime := date.Add(collectDataEndTimeOfDay)

	datacollection.GetTransitDataWithGoogleAPI(startTime, endTime, collectDataInterval, collectDataAPIKeyFile)
}
