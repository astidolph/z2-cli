package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/strava-cli/internal/auth"
	"github.com/strava-cli/internal/strava"
)

var (
	weeksBack   int
	dayFilter   string
	showAll     bool
	minDistance  float64
)

var runsCmd = &cobra.Command{
	Use:   "runs",
	Short: "Display your running data",
	Long: `Fetch and display your runs from Strava.

By default, only shows zone 2 runs (requires zone 2 HR to be set via 'strava-cli config').
Use --all to show all runs regardless of heart rate.
Use --day to filter to a specific day of the week.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := getValidToken()
		if err != nil {
			return err
		}

		client := strava.NewClient(token.AccessToken)

		since := time.Now().AddDate(0, 0, -weeksBack*7)
		runs, err := client.GetAllRunsSince(since)
		if err != nil {
			return fmt.Errorf("could not fetch runs: %w", err)
		}

		if dayFilter != "" {
			day, err := parseWeekday(dayFilter)
			if err != nil {
				return err
			}
			runs = strava.FilterByWeekday(runs, day)
		}

		if minDistance > 0 {
			runs = strava.FilterByMinDistance(runs, minDistance)
		}

		if !showAll {
			config, err := auth.LoadConfig()
			if err != nil {
				return err
			}
			if config.Zone2HR == 0 {
				return fmt.Errorf("zone 2 HR not set — run 'strava-cli config --zone2-hr <value>' or use --all to skip filtering")
			}
			runs = strava.FilterByMaxHR(runs, float64(config.Zone2HR))
			fmt.Printf("Zone 2 runs (avg HR ≤ %d bpm) from the last %d weeks:\n\n", config.Zone2HR, weeksBack)
		} else {
			fmt.Printf("All runs from the last %d weeks:\n\n", weeksBack)
		}

		if len(runs) == 0 {
			fmt.Println("No matching runs found.")
			return nil
		}

		printRunsTable(runs)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runsCmd)
	runsCmd.Flags().IntVarP(&weeksBack, "weeks", "w", 12, "Number of weeks to look back")
	runsCmd.Flags().StringVarP(&dayFilter, "day", "d", "", "Day of week to filter (e.g. sunday, monday)")
	runsCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all runs, skip zone 2 filtering")
	runsCmd.Flags().Float64Var(&minDistance, "min-distance", 0, "Minimum distance in km (e.g. 12 for long runs)")
}

func getValidToken() (*auth.Token, error) {
	config, err := auth.LoadConfig()
	if err != nil {
		return nil, err
	}

	token, err := auth.LoadToken()
	if err != nil {
		return nil, err
	}

	if token.IsExpired() {
		token, err = auth.RefreshAccessToken(config.ClientID, config.ClientSecret, token)
		if err != nil {
			return nil, fmt.Errorf("could not refresh token: %w", err)
		}
		if err := auth.SaveToken(token); err != nil {
			return nil, fmt.Errorf("could not save refreshed token: %w", err)
		}
	}

	return token, nil
}

func parseWeekday(s string) (time.Weekday, error) {
	days := map[string]time.Weekday{
		"sunday":    time.Sunday,
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
	}
	day, ok := days[s]
	if !ok {
		return 0, fmt.Errorf("invalid day: %s", s)
	}
	return day, nil
}

func printRunsTable(runs []strava.Activity) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DATE\tDISTANCE (km)\tTIME\tAVG HR\tPACE (/km)")
	fmt.Fprintln(w, "────\t─────────────\t────\t──────\t──────────")

	for _, r := range runs {
		t, _ := r.StartTime()
		date := t.Format("02 Jan 2006")
		distKm := r.Distance / 1000.0
		duration := formatDuration(r.MovingTime)
		pace := formatPace(r.Distance, r.MovingTime)

		hr := "—"
		if r.HasHeartrate {
			hr = fmt.Sprintf("%.0f bpm", r.AverageHeartrate)
		}

		fmt.Fprintf(w, "%s\t%.2f\t%s\t%s\t%s\n", date, distKm, duration, hr, pace)
	}
	w.Flush()
}

func formatDuration(seconds int) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	if h > 0 {
		return fmt.Sprintf("%dh %02dm %02ds", h, m, s)
	}
	return fmt.Sprintf("%dm %02ds", m, s)
}

func formatPace(distMeters float64, seconds int) string {
	if distMeters == 0 {
		return "—"
	}
	paceSeconds := float64(seconds) / (distMeters / 1000.0)
	m := int(paceSeconds) / 60
	s := int(paceSeconds) % 60
	return fmt.Sprintf("%d:%02d", m, s)
}
