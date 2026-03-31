package chart

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/z2-cli/internal/stats"
	"github.com/z2-cli/internal/strava"
)

const kmToMile = 1.60934

type ChartData struct {
	Dates      []string
	EF         []opts.LineData
	Pace       []opts.LineData
	PaceMi     []opts.LineData
	Distance   []opts.LineData
	DistanceMi []opts.LineData
	HR         []opts.LineData
}

func BuildChartData(runs []strava.Activity) ChartData {
	// Sort by date ascending for charts
	sorted := make([]strava.Activity, len(runs))
	copy(sorted, runs)
	sortByDateAsc(sorted)

	data := ChartData{}
	for _, r := range sorted {
		t, err := r.StartTime()
		if err != nil {
			continue
		}
		data.Dates = append(data.Dates, t.Format("02 Jan"))

		ef := stats.EfficiencyFactor(r)
		data.EF = append(data.EF, opts.LineData{Value: fmt.Sprintf("%.4f", ef)})

		if r.Distance > 0 && r.MovingTime > 0 {
			paceKm := float64(r.MovingTime) / (r.Distance / 1000.0) / 60.0
			paceMi := paceKm * kmToMile
			data.Pace = append(data.Pace, opts.LineData{Value: fmt.Sprintf("%.2f", paceKm)})
			data.PaceMi = append(data.PaceMi, opts.LineData{Value: fmt.Sprintf("%.2f", paceMi)})
		} else {
			data.Pace = append(data.Pace, opts.LineData{Value: nil})
			data.PaceMi = append(data.PaceMi, opts.LineData{Value: nil})
		}

		distKm := r.Distance / 1000.0
		distMi := distKm / kmToMile
		data.Distance = append(data.Distance, opts.LineData{Value: fmt.Sprintf("%.2f", distKm)})
		data.DistanceMi = append(data.DistanceMi, opts.LineData{Value: fmt.Sprintf("%.2f", distMi)})

		if r.HasHeartrate {
			data.HR = append(data.HR, opts.LineData{Value: fmt.Sprintf("%.0f", r.AverageHeartrate)})
		} else {
			data.HR = append(data.HR, opts.LineData{Value: nil})
		}
	}
	return data
}

func dataRange(series []opts.LineData, padding float64) (float64, float64) {
	min := math.MaxFloat64
	max := -math.MaxFloat64
	for _, d := range series {
		var v float64
		switch val := d.Value.(type) {
		case string:
			v, _ = strconv.ParseFloat(val, 64)
		case float64:
			v = val
		case int:
			v = float64(val)
		default:
			continue
		}
		if v == 0 {
			continue
		}
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	spread := max - min
	return min - spread*padding, max + spread*padding
}

func RenderEF(w io.Writer, data ChartData) error {
	efMin, efMax := dataRange(data.EF, 0.2)
	hrMin, hrMax := dataRange(data.HR, 0.2)

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Efficiency Factor & Heart Rate Over Time",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name:      "EF (speed/HR)",
			Min:       fmt.Sprintf("%.4f", efMin),
			Max:       fmt.Sprintf("%.4f", efMax),
			AxisLabel: &opts.AxisLabel{Formatter: "{value}"},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    opts.Bool(true),
			Trigger: "axis",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "500px",
			Theme:  "dark",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(true),
		}),
	)

	line.ExtendYAxis(opts.YAxis{
		Name:      "HR (bpm)",
		Min:       math.Floor(hrMin),
		Max:       math.Ceil(hrMax),
		AxisLabel: &opts.AxisLabel{Formatter: "{value}"},
	})

	line.SetXAxis(data.Dates).
		AddSeries("EF", data.EF).
		AddSeries("Avg HR", data.HR, charts.WithLineChartOpts(opts.LineChart{
			YAxisIndex: 1,
		}))

	return line.Render(w)
}

func RenderPace(w io.Writer, data ChartData) error {
	line := newLine("Pace Over Time", "Pace (min)")
	line.SetXAxis(data.Dates).
		AddSeries("Pace /km", data.Pace).
		AddSeries("Pace /mi", data.PaceMi)
	return line.Render(w)
}

func RenderDistance(w io.Writer, data ChartData) error {
	line := newLine("Distance Over Time", "Distance")
	line.SetXAxis(data.Dates).
		AddSeries("km", data.Distance).
		AddSeries("mi", data.DistanceMi)
	return line.Render(w)
}

func RenderHR(w io.Writer, data ChartData) error {
	line := newLine("Average Heart Rate Over Time", "HR (bpm)")
	line.SetXAxis(data.Dates).
		AddSeries("Avg HR", data.HR)
	return line.Render(w)
}

func newLine(title, yAxisName string) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: title,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: yAxisName,
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    opts.Bool(true),
			Trigger: "axis",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "500px",
			Theme:  "dark",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(true),
		}),
	)
	return line
}

func sortByDateAsc(runs []strava.Activity) {
	sort.Slice(runs, func(i, j int) bool {
		ti, _ := runs[i].StartTime()
		tj, _ := runs[j].StartTime()
		return ti.Before(tj)
	})
}

// RenderAll generates a single HTML page with all charts.
func RenderAll(w io.Writer, data ChartData) error {
	page := `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>z2-cli — Zone 2 Training Charts</title>
<style>body{background:#100c2a;margin:0;padding:20px;font-family:sans-serif}
.chart{margin-bottom:20px}</style></head><body>`

	if _, err := fmt.Fprint(w, page); err != nil {
		return err
	}

	renderers := []struct {
		name string
		fn   func(io.Writer, ChartData) error
	}{
		{"ef", RenderEF},
		{"pace", RenderPace},
		{"distance", RenderDistance},
		{"hr", RenderHR},
	}

	for _, r := range renderers {
		if _, err := fmt.Fprint(w, `<div class="chart">`); err != nil {
			return err
		}
		if err := r.fn(w, data); err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, `</div>`); err != nil {
			return err
		}
	}

	_, err := fmt.Fprint(w, `</body></html>`)
	return err
}

// AvailableTypes returns the valid chart type names.
func AvailableTypes() []string {
	return []string{"ef", "pace", "distance", "hr", "all"}
}

// RenderByType renders a chart by type name.
func RenderByType(w io.Writer, data ChartData, chartType string) error {
	switch chartType {
	case "ef":
		return RenderEF(w, data)
	case "pace":
		return RenderPace(w, data)
	case "distance":
		return RenderDistance(w, data)
	case "hr":
		return RenderHR(w, data)
	case "all":
		return RenderAll(w, data)
	default:
		return fmt.Errorf("unknown chart type: %s (options: ef, pace, distance, hr, all)", chartType)
	}
}
