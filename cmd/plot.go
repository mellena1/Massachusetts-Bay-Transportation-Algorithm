package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/calculation"
	"github.com/mellena1/Massachusetts-Bay-Transportation-Algorithm/datacollection"
	"github.com/spf13/cobra"
)

var plotInputFile string
var plotOutputFile string
var plotStopA string
var plotStopB string

func init() {
	plotCmd.Flags().StringVarP(&plotInputFile, "input", "i", "", "File read the edge data from (required)")
	plotCmd.MarkFlagRequired("input")

	plotCmd.Flags().StringVarP(&plotOutputFile, "output", "o", "", "File to save the plot to (required)")
	plotCmd.MarkFlagRequired("output")

	plotCmd.Flags().StringVar(&plotStopA, "stopA", "", "First stop of the edge to plot")
	plotCmd.Flags().StringVar(&plotStopB, "stopB", "", "Last stop of the edge to plot")

	rootCmd.AddCommand(plotCmd)
}

var plotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Given a data file, plot the cubic spline approximation into a graph image",
	Long:  "Given a data file, plot the cubic spline approximation into a graph image",
	Run:   plotFunc,
}

func plotFunc(cmd *cobra.Command, args []string) {
	edgeData, err := datacollection.ImportEdgeData(plotInputFile)
	if err != nil {
		log.Fatalf("failed reading in edge data: %s", err)
	}

	plotStopA = makeStopNameValid(plotStopA)
	plotStopB = makeStopNameValid(plotStopB)

	if plotStopA == "" && plotStopB == "" {
		calculation.PlotAllEdges(edgeData, plotOutputFile)
		return
	}

	if (plotStopA != "" && plotStopB == "") || (plotStopA == "" && plotStopB != "") {
		fmt.Println("Please specify both stopA and stopB")
		os.Exit(1)
	}

	edgeKey := datacollection.GetEdgeKey(plotStopA, plotStopB)
	if _, ok := edgeData[edgeKey]; !ok {
		fmt.Println("Invalid edge. Please give a valid stopA and stopB name")
		os.Exit(1)
	}
	calculation.PlotEdge(edgeData, edgeKey, plotOutputFile)

	walkingEdgeKey := datacollection.GetEdgeKeyWalking(plotStopA, plotStopB)
	if _, ok := edgeData[walkingEdgeKey]; ok {
		calculation.PlotEdge(edgeData, walkingEdgeKey, plotOutputFile)
	}
}
