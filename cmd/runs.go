package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/strava-cli/internal/auth"
	"github.com/strava-cli/internal/stats"
	"github.com/strava-cli/internal/strava"
)

var (
	weeksBack   int
	dayFilter   string
	showAll     bool
	minDistance  float64
	sortBy      string
	ascending   bool
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

		now := time.Now()
		since := now.AddDate(0, 0, -weeksBack*7)
		priorSince := since.AddDate(0, 0, -weeksBack*7)

		runs, err := client.GetAllRunsSince(priorSince)
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

		var currentRuns, priorRuns []strava.Activity
		for _, r := range runs {
			t, err := r.StartTime()
			if err != nil {
				continue
			}
			if t.After(since) {
				currentRuns = append(currentRuns, r)
			} else {
				priorRuns = append(priorRuns, r)
			}
		}

		if len(currentRuns) == 0 {
			fmt.Println("No matching runs found.")
			return nil
		}

		if err := sortRuns(currentRuns, sortBy, ascending); err != nil {
			return err
		}

		printRunsTable(currentRuns)
		printSummary(currentRuns, priorRuns, weeksBack)
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

func sortRuns(runs []strava.Activity, by string, asc bool) error {
	var less func(i, j int) bool

	switch by {
	case "date":
		less = func(i, j int) bool {
			ti, _ := runs[i].StartTime()
			tj, _ := runs[j].StartTime()
			return ti.After(tj)
		}
	case "distance":
		less = func(i, j int) bool {
			return runs[i].Distance > runs[j].Distance
		}
	case "time":
		less = func(i, j int) bool {
			return runs[i].MovingTime > runs[j].MovingTime
		}
	case "hr":
		less = func(i, j int) bool {
			return runs[i].AverageHeartrate > runs[j].AverageHeartrate
		}
	case "pace":
		less = func(i, j int) bool {
			// Lower pace seconds = faster, so "descending" means fastest first
			pi := paceSecondsPerKm(runs[i])
			pj := paceSecondsPerKm(runs[j])
			return pi < pj
		}
	case "ef":
		less = func(i, j int) bool {
			return stats.EfficiencyFactor(runs[i]) > stats.EfficiencyFactor(runs[j])
		}
	default:
		return fmt.Errorf("invalid sort column: %s (options: date, distance, time, hr, pace, ef)", by)
	}

	if asc {
		original := less
		less = func(i, j int) bool { return !original(i, j) }
	}

	sort.SliceStable(runs, less)
	return nil
}

func paceSecondsPerKm(a strava.Activity) float64 {
	if a.Distance == 0 {
		return 0
	}
	return float64(a.MovingTime) / (a.Distance / 1000.0)
}

func printRunsTable(runs []strava.Activity) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DATE\tDISTANCE (km)\tTIME\tAVG HR\tPACE (/km)\tPACE (/mi)\tEF")
	fmt.Fprintln(w, "────\t─────────────\t────\t──────\t──────────\t──────────\t──")

	for _, r := range runs {
		t, _ := r.StartTime()
		date := t.Format("02 Jan 2006")
		distKm := r.Distance / 1000.0
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

		fmt.Fprintf(w, "%s\t%.2f\t%s\t%s\t%s\t%s\t%s\n", date, distKm, duration, hr, paceKm, paceMi, efStr)
	}
	w.Flush()
}

func printSummary(current, prior []strava.Activity, weeks int) {
	cur := stats.Summarise(current)

	fmt.Printf("\nSummary (last %d weeks, %d runs, %.1f km total):\n", weeks, cur.Count, cur.TotalKm)

	if cur.AvgEF > 0 {
		efLine := fmt.Sprintf("  Avg EF:   %.4f", cur.AvgEF)
		if len(prior) > 0 {
			prev := stats.Summarise(prior)
			if prev.AvgEF > 0 {
				trend := stats.TrendPercent(cur, prev)
				arrow := "→"
				if trend > 0 {
					arrow = "↑"
				} else if trend < 0 {
					arrow = "↓"
				}
				efLine += fmt.Sprintf(" %s (%+.1f%% vs prior %d weeks)", arrow, trend, weeks)
			}
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
