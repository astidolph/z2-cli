package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "strava-cli",
	Short: "A CLI tool to track your Strava zone 2 training progress",
	Long:  "Fetch and visualise your Sunday long run data from Strava — distance, average HR, and time.",
}

func Execute() error {
	return rootCmd.Execute()
}
