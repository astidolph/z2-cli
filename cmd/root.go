package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/z2-cli/internal/storage"
)

var rootCmd = &cobra.Command{
	Use:   "z2-cli",
	Short: "Track your zone 2 training progress from Strava",
	Long:  "Fetch and visualise your zone 2 running data from Strava — distance, heart rate, pace, and efficiency factor.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
			store, err := storage.NewPGStore(context.Background(), dbURL)
			if err != nil {
				return fmt.Errorf("could not connect to database: %w", err)
			}
			storage.Init(store)
		} else {
			store, err := storage.NewFileStore()
			if err != nil {
				return fmt.Errorf("could not initialize file storage: %w", err)
			}
			storage.Init(store)
		}
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
