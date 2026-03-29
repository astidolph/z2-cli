package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/z2-cli/internal/chart"
	"github.com/z2-cli/internal/service"
)

var (
	chartType      string
	chartWeeksBack int
	chartDayFilter string
	chartShowAll   bool
	chartMinDist   float64
)

var chartCmd = &cobra.Command{
	Use:   "chart",
	Short: "Generate training charts in your browser",
	Long: `Generate interactive HTML charts from your run data and open them in the browser.

Chart types:
  ef        Efficiency factor over time (default)
  pace      Pace per km and per mile over time
  distance  Distance over time in km and miles
  hr        Average heart rate over time
  all       All charts on a single page`,
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := service.FetchRuns(service.RunsQuery{
			WeeksBack:  chartWeeksBack,
			Day:        chartDayFilter,
			MinDistance: chartMinDist,
			ShowAll:    chartShowAll,
			SortBy:     "date",
			Ascending:  true,
		})
		if err != nil {
			return err
		}

		if len(result.CurrentRuns) == 0 {
			fmt.Println("No matching runs found.")
			return nil
		}

		data := chart.BuildChartData(result.CurrentRuns)

		tmpDir := os.TempDir()
		filePath := filepath.Join(tmpDir, "z2-cli-chart.html")

		f, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("could not create chart file: %w", err)
		}
		defer f.Close()

		if err := chart.RenderByType(f, data, chartType); err != nil {
			return err
		}

		fmt.Printf("Chart generated: %s\n", filePath)
		return openBrowser(filePath)
	},
}

func init() {
	rootCmd.AddCommand(chartCmd)
	chartCmd.Flags().StringVarP(&chartType, "type", "t", "ef", "Chart type: ef, pace, distance, hr, all")
	chartCmd.Flags().IntVarP(&chartWeeksBack, "weeks", "w", 12, "Number of weeks to look back")
	chartCmd.Flags().StringVarP(&chartDayFilter, "day", "d", "", "Day of week to filter")
	chartCmd.Flags().BoolVarP(&chartShowAll, "all", "a", false, "Show all runs, skip zone 2 filtering")
	chartCmd.Flags().Float64Var(&chartMinDist, "min-distance", 0, "Minimum distance in km")
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
