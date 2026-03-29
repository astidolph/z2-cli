package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "z2-cli",
	Short: "Track your zone 2 training progress from Strava",
	Long:  "Fetch and visualise your zone 2 running data from Strava — distance, heart rate, pace, and efficiency factor.",
}

func Execute() error {
	return rootCmd.Execute()
}
