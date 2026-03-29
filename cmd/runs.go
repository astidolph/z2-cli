package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/z2-cli/internal/service"
	"github.com/z2-cli/internal/stats"
	"github.com/z2-cli/internal/strava"
)

var (
	weeksBack  int
	dayFilter  string
	showAll    bool
	minDistance float64
	sortBy     string
	ascending  bool
)

var runsCmd = &cobra.Command{
	Use:   "runs",
	Short: "Display your running data",
	Long: `Fetch and display your runs from Strava.

By default, only shows zone 2 runs (requires zone 2 HR to be set via 'z2-cli config').
Use --all to show all runs regardless of heart rate.
Use --day to filter to a specific day of the week.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := service.FetchRuns(service.RunsQuery{
			WeeksBack:  weeksBack,
			Day:        dayFilter,
			MinDistance: minDistance,
			ShowAll:    showAll,
			SortBy:     sortBy,
			Ascending:  ascending,
		})
		if err != nil {
			return err
		}

		if result.Zone2HR > 0 {
			fmt.Printf("Zone 2 runs (avg HR ≤ %d bpm) from the last %d weeks:\n\n", result.Zone2HR, weeksBack)
		} else {
			fmt.Printf("All runs from the last %d weeks:\n\n", weeksBack)
		}

		if len(result.CurrentRuns) == 0 {
			fmt.Println("No matching runs found.")
			return nil
		}

		printRunsTable(result.CurrentRuns)
		printSummary(result)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runsCmd)
	runsCmd.Flags().IntVarP(&weeksBack, "weeks", "w", 12, "Number of weeks to look back")
	runsCmd.Flags().StringVarP(&dayFilter, "day", "d", "", "Day of week to filter (e.g. sunday, monday)")
	runsCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all runs, skip zone 2 filtering")
	runsCmd.Flags().Float64Var(&minDistance, "min-distance", 0, "Minimum distance in km (e.g. 12 for long runs)")
	runsCmd.Flags().StringVar(&sortBy, "sort", "date", "Sort by: date, distance, time, hr, pace, ef")
	runsCmd.Flags().BoolVar(&ascending, "asc", false, "Sort in ascending order (default is descending)")
}

func printRunsTable(runs []strava.Activity) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DATE\tDIST (km)\tDIST (mi)\tTIME\tAVG HR\tPACE (/km)\tPACE (/mi)\tEF")
	fmt.Fprintln(w, "────\t─────────\t─────────\t────\t──────\t──────────\t──────────\t──")

	for _, r := range runs {
		t, _ := r.StartTime()
		date := t.Format("02 Jan 2006")
		distKm := r.Distance / 1000.0
		distMi := distKm / kmToMile
		duration := formatDuration(r.MovingTime)
		paceKm := formatPacePerKm(r.Distance, r.MovingTime)
		paceMi := formatPacePerMile(r.Distance, r.MovingTime)

		hr := "—"
		if r.HasHeartrate {
			hr = fmt.Sprintf("%.0f bpm", r.AverageHeartrate)
		}

		ef := stats.EfficiencyFactor(r)
		efStr := "—"
		if ef > 0 {
			efStr = fmt.Sprintf("%.4f", ef)
		}

		fmt.Fprintf(w, "%s\t%.2f\t%.2f\t%s\t%s\t%s\t%s\t%s\n", date, distKm, distMi, duration, hr, paceKm, paceMi, efStr)
	}
	w.Flush()
}

func printSummary(result *service.RunsResult) {
	cur := result.Current

	totalMi := cur.TotalKm / kmToMile
	fmt.Printf("\nSummary (last %d weeks, %d runs, %.1f km / %.1f mi total):\n", result.WeeksBack, cur.Count, cur.TotalKm, totalMi)

	if cur.AvgEF > 0 {
		efLine := fmt.Sprintf("  Avg EF:   %.4f", cur.AvgEF)
		if result.Prior.AvgEF > 0 {
			trend := stats.TrendPercent(cur, result.Prior)
			arrow := "→"
			if trend > 0 {
				arrow = "↑"
			} else if trend < 0 {
				arrow = "↓"
			}
			efLine += fmt.Sprintf(" %s (%+.1f%% vs prior %d weeks)", arrow, trend, result.WeeksBack)
		}
		fmt.Println(efLine)
	}

	if cur.AvgHR > 0 {
		fmt.Printf("  Avg HR:   %.0f bpm\n", cur.AvgHR)
	}
	if cur.AvgPace > 0 {
		pacePerMile := cur.AvgPace * kmToMile
		fmt.Printf("  Avg Pace: %s/km (%s/mi)\n", formatPaceSeconds(cur.AvgPace), formatPaceSeconds(pacePerMile))
	}
}

func formatPaceSeconds(totalSeconds float64) string {
	m := int(totalSeconds) / 60
	s := int(totalSeconds) % 60
	return fmt.Sprintf("%d:%02d", m, s)
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

const kmToMile = 1.60934

func formatPacePerKm(distMeters float64, seconds int) string {
	if distMeters == 0 {
		return "—"
	}
	paceSeconds := float64(seconds) / (distMeters / 1000.0)
	m := int(paceSeconds) / 60
	s := int(paceSeconds) % 60
	return fmt.Sprintf("%d:%02d", m, s)
}

func formatPacePerMile(distMeters float64, seconds int) string {
	if distMeters == 0 {
		return "—"
	}
	distMiles := distMeters / 1000.0 / kmToMile
	paceSeconds := float64(seconds) / distMiles
	m := int(paceSeconds) / 60
	s := int(paceSeconds) % 60
	return fmt.Sprintf("%d:%02d", m, s)
}
