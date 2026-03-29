package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/strava-cli/internal/auth"
)

var (
	zone2HR int
	age     int
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure your training settings",
	Long: `Set your zone 2 heart rate ceiling to enable filtering.

You can set it directly or calculate from your age using the Maffetone formula (180 - age).

Examples:
  strava-cli config --zone2-hr 150
  strava-cli config --age 30          # sets zone 2 HR to 150 (180 - 30)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		hrChanged := cmd.Flags().Changed("zone2-hr")
		ageChanged := cmd.Flags().Changed("age")

		if !hrChanged && !ageChanged {
			config, err := auth.LoadConfig()
			if err != nil {
				fmt.Println("No configuration found. Use --zone2-hr or --age to set your zone 2 ceiling.")
				return nil
			}
			if config.Zone2HR == 0 {
				fmt.Println("Zone 2 HR not set. Use --zone2-hr or --age to set your zone 2 ceiling.")
			} else {
				fmt.Printf("Zone 2 HR ceiling: %d bpm\n", config.Zone2HR)
				fmt.Println("Runs with average HR at or below this value are considered zone 2.")
			}
			return nil
		}

		if hrChanged && ageChanged {
			return fmt.Errorf("use either --zone2-hr or --age, not both")
		}

		hr := zone2HR
		if ageChanged {
			if age < 10 || age > 100 {
				return fmt.Errorf("age should be between 10 and 100, got %d", age)
			}
			hr = 180 - age
		}

		if hr < 80 || hr > 200 {
			return fmt.Errorf("zone 2 HR should be between 80 and 200, got %d", hr)
		}

		config, err := auth.LoadConfig()
		if err != nil {
			config = &auth.Config{}
		}

		config.Zone2HR = hr
		if err := auth.SaveConfig(config); err != nil {
			return fmt.Errorf("could not save config: %w", err)
		}

		if ageChanged {
			fmt.Printf("Zone 2 HR ceiling set to %d bpm (180 - %d)\n", hr, age)
		} else {
			fmt.Printf("Zone 2 HR ceiling set to %d bpm\n", hr)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().IntVar(&zone2HR, "zone2-hr", 0, "Zone 2 heart rate ceiling in bpm")
	configCmd.Flags().IntVar(&age, "age", 0, "Your age (calculates zone 2 HR as 180 - age)")
}
