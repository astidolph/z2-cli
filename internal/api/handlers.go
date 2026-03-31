package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/z2-cli/internal/auth"
	"github.com/z2-cli/internal/cache"
	"github.com/z2-cli/internal/chart"
	"github.com/z2-cli/internal/service"
	"github.com/z2-cli/internal/stats"
)

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	token, err := auth.LoadToken()
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"authenticated": false,
			"message":       "Not authenticated — run 'z2-cli auth' to connect Strava",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"authenticated": !token.IsExpired(),
		"expires_at":    token.ExpiresAt,
	})
}

func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	config, err := auth.LoadConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"zone2_hr": config.Zone2HR,
	})
}

func handlePutConfig(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Zone2HR *int `json:"zone2_hr"`
		Age     *int `json:"age"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	config, err := auth.LoadConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	switch {
	case body.Zone2HR != nil:
		if *body.Zone2HR <= 0 {
			writeError(w, http.StatusBadRequest, "zone2_hr must be positive")
			return
		}
		config.Zone2HR = *body.Zone2HR
	case body.Age != nil:
		if *body.Age <= 0 {
			writeError(w, http.StatusBadRequest, "age must be positive")
			return
		}
		config.Zone2HR = 180 - *body.Age
	default:
		writeError(w, http.StatusBadRequest, "provide zone2_hr or age")
		return
	}

	if err := auth.SaveConfig(config); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"zone2_hr": config.Zone2HR,
	})
}

func handleGetRuns(w http.ResponseWriter, r *http.Request) {
	query, err := parseRunsQuery(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := service.FetchRuns(query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Add EF and trend to the response
	type runsResponse struct {
		*service.RunsResult
		EFTrend float64 `json:"ef_trend"`
	}

	writeJSON(w, http.StatusOK, runsResponse{
		RunsResult: result,
		EFTrend:    stats.TrendPercent(result.Current, result.Prior),
	})
}

type chartDataResponse struct {
	Dates      []string   `json:"dates"`
	EF         []*float64 `json:"ef"`
	Pace       []*float64 `json:"pace"`
	PaceMi     []*float64 `json:"pace_mi"`
	Distance   []*float64 `json:"distance"`
	DistanceMi []*float64 `json:"distance_mi"`
	HR         []*float64 `json:"hr"`
}

func handleGetChartData(w http.ResponseWriter, r *http.Request) {
	query, err := parseRunsQuery(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := service.FetchRuns(query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := chart.BuildChartData(result.CurrentRuns)
	resp := chartDataResponse{
		Dates:      data.Dates,
		EF:         lineDataToFloats(data.EF),
		Pace:       lineDataToFloats(data.Pace),
		PaceMi:     lineDataToFloats(data.PaceMi),
		Distance:   lineDataToFloats(data.Distance),
		DistanceMi: lineDataToFloats(data.DistanceMi),
		HR:         lineDataToFloats(data.HR),
	}

	writeJSON(w, http.StatusOK, resp)
}

func handleRefresh(w http.ResponseWriter, r *http.Request) {
	if err := cache.Invalidate(); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "cache cleared"})
}

// parseRunsQuery extracts RunsQuery fields from URL query parameters.
func parseRunsQuery(r *http.Request) (service.RunsQuery, error) {
	q := r.URL.Query()
	query := service.RunsQuery{
		WeeksBack: 12,
		SortBy:    "date",
	}

	if v := q.Get("weeks"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			return query, fmt.Errorf("invalid weeks parameter")
		}
		query.WeeksBack = n
	}
	if v := q.Get("day"); v != "" {
		query.Day = v
	}
	if v := q.Get("minDistance"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil || f < 0 {
			return query, fmt.Errorf("invalid minDistance parameter")
		}
		query.MinDistance = f
	}
	if v := q.Get("all"); v == "true" {
		query.ShowAll = true
	}
	if v := q.Get("sort"); v != "" {
		query.SortBy = v
	}
	if v := q.Get("asc"); v == "true" {
		query.Ascending = true
	}
	if v := q.Get("refresh"); v == "true" {
		query.ForceRefresh = true
	}

	return query, nil
}

// lineDataToFloats converts go-echarts LineData values to plain float64 pointers.
// nil Values are preserved as nil so they serialize to JSON null (Chart.js skips nulls).
func lineDataToFloats(data []opts.LineData) []*float64 {
	out := make([]*float64, len(data))
	for i, d := range data {
		if d.Value == nil {
			continue
		}
		var f float64
		switch v := d.Value.(type) {
		case string:
			parsed, err := strconv.ParseFloat(v, 64)
			if err != nil {
				continue
			}
			f = parsed
		case float64:
			f = v
		case int:
			f = float64(v)
		default:
			continue
		}
		out[i] = &f
	}
	return out
}
