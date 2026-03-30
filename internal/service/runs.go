package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/z2-cli/internal/auth"
	"github.com/z2-cli/internal/cache"
	"github.com/z2-cli/internal/stats"
	"github.com/z2-cli/internal/strava"
)

type RunsQuery struct {
	WeeksBack    int
	Day          string
	MinDistance   float64
	ShowAll      bool
	SortBy       string
	Ascending    bool
	ForceRefresh bool
}

type RunsResult struct {
	CurrentRuns []strava.Activity `json:"current_runs"`
	PriorRuns   []strava.Activity `json:"prior_runs"`
	Current     stats.Summary     `json:"current"`
	Prior       stats.Summary     `json:"prior"`
	Zone2HR     int               `json:"zone2_hr"`
	WeeksBack   int               `json:"weeks_back"`
}

func FetchRuns(query RunsQuery) (*RunsResult, error) {
	now := time.Now()
	since := now.AddDate(0, 0, -query.WeeksBack*7)
	priorSince := since.AddDate(0, 0, -query.WeeksBack*7)

	runs, err := fetchActivities(priorSince, query.ForceRefresh)
	if err != nil {
		return nil, err
	}

	if query.Day != "" {
		day, err := parseWeekday(query.Day)
		if err != nil {
			return nil, err
		}
		runs = strava.FilterByWeekday(runs, day)
	}

	if query.MinDistance > 0 {
		runs = strava.FilterByMinDistance(runs, query.MinDistance)
	}

	var zone2HR int
	if !query.ShowAll {
		config, err := auth.LoadConfig()
		if err != nil {
			return nil, err
		}
		if config.Zone2HR == 0 {
			return nil, fmt.Errorf("zone 2 HR not set — run 'z2-cli config --zone2-hr <value>' or use --all to skip filtering")
		}
		zone2HR = config.Zone2HR
		runs = strava.FilterByMaxHR(runs, float64(config.Zone2HR))
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

	if err := SortRuns(currentRuns, query.SortBy, query.Ascending); err != nil {
		return nil, err
	}

	return &RunsResult{
		CurrentRuns: currentRuns,
		PriorRuns:   priorRuns,
		Current:     stats.Summarise(currentRuns),
		Prior:       stats.Summarise(priorRuns),
		Zone2HR:     zone2HR,
		WeeksBack:   query.WeeksBack,
	}, nil
}

func fetchActivities(since time.Time, forceRefresh bool) ([]strava.Activity, error) {
	if !forceRefresh {
		if cached := cache.Load(); cached != nil && cached.IsFresh(since) {
			return cached.Activities, nil
		}
	}

	token, err := GetValidToken()
	if err != nil {
		return nil, err
	}

	client := strava.NewClient(token.AccessToken)
	runs, err := client.GetAllRunsSince(since)
	if err != nil {
		return nil, fmt.Errorf("could not fetch runs: %w", err)
	}

	_ = cache.Save(&cache.CachedData{
		FetchedAt:  time.Now(),
		SinceUnix:  since.Unix(),
		Activities: runs,
	})

	return runs, nil
}

func GetValidToken() (*auth.Token, error) {
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

func SortRuns(runs []strava.Activity, by string, asc bool) error {
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
