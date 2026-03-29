package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/strava-cli/internal/auth"
)

var maxHR int

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure your training settings",
	Long: `Set your max heart rate to enable zone 2 filtering.

Zone 2 is calculated as 60-70% of your max heart rate.
Example: max HR of 190 → zone 2 is 114-133 bpm.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed("max-hr") {
			config, err := auth.LoadConfig()
			if err != nil {
				fmt.Println("No configuration found. Use --max-hr to set your max heart rate.")
				return nil
			}
			if config.MaxHR == 0 {
				fmt.Println("Max HR not set. Use --max-hr to set your max heart rate.")
			} else {
				low, high := config.Zone2Range()
				fmt.Printf("Max HR:  %d bpm\n", config.MaxHR)
				fmt.Printf("Zone 2:  %.0f-%.0f bpm\n", low, high)
			}
			return nil
		}

		if maxHR < 100 || maxHR > 250 {
			return fmt.Errorf("max HR should be between 100 and 250, got %d", maxHR)
		}

		config, err := auth.LoadConfig()
		if err != nil {
			config = &auth.Config{}
		}

		config.MaxHR = maxHR
		if err := auth.SaveConfig(config); err != nil {
			return fmt.Errorf("could not save config: %w", err)
		}

		low, high := config.Zone2Range()
		fmt.Printf("Max HR set to %d bpm. Zone 2 range: %.0f-%.0f bpm\n", maxHR, low, high)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().IntVar(&maxHR, "max-hr", 0, "Your maximum heart rate in bpm")
}
