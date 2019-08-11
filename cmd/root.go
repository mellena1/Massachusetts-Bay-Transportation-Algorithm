package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mbta",
	Short: "mbta finds the fastest path to hit every mbta stop",
	Long:  "mbta allows you to download data from the google maps api and use that data to calculate",
}

// Execute runs the cli tool
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
