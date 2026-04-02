package service

import (
	"fmt"
	"time"

	"github.com/z2-cli/internal/cache"
	"github.com/z2-cli/internal/stats"
	"github.com/z2-cli/internal/strava"
)

const LeaderboardPageSize = 30

type LeaderboardResult struct {
	Runs       []strava.Activity `json:"runs"`
	TotalCount int               `json:"total_count"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
}

// FetchLeaderboard returns a page of runs sorted by EF descending.
// Only runs with heart rate data (EF > 0) are included.
func FetchLeaderboard(page int) (*LeaderboardResult, error) {
	history := cache.LoadHistory()
	if history == nil || len(history.Activities) == 0 {
		return &LeaderboardResult{
			Runs:       []strava.Activity{},
			TotalCount: 0,
			Page:       page,
			PageSize:   LeaderboardPageSize,
		}, nil
	}

	// Filter to runs with HR data, minimum 3 miles, and sort by EF descending.
	const minDistanceMeters = 4828.03 // 3 miles
	var eligible []strava.Activity
	for _, a := range history.Activities {
		if a.HasHeartrate && a.Distance >= minDistanceMeters && stats.EfficiencyFactor(a) > 0 {
			eligible = append(eligible, a)
		}
	}

	if err := SortRuns(eligible, "ef", false); err != nil {
		return nil, err
	}

	total := len(eligible)

	// Paginate.
	start := (page - 1) * LeaderboardPageSize
	if start >= total {
		return &LeaderboardResult{
			Runs:       []strava.Activity{},
			TotalCount: total,
			Page:       page,
			PageSize:   LeaderboardPageSize,
		}, nil
	}
	end := start + LeaderboardPageSize
	if end > total {
		end = total
	}

	return &LeaderboardResult{
		Runs:       eligible[start:end],
		TotalCount: total,
		Page:       page,
		PageSize:   LeaderboardPageSize,
	}, nil
}

// RefreshLeaderboard performs an incremental sync of the full run history
// from Strava. If no history exists, it fetches all runs from the beginning.
func RefreshLeaderboard() error {
	token, err := GetValidToken()
	if err != nil {
		return err
	}

	client := strava.NewClient(token.AccessToken)

	history := cache.LoadHistory()
	var since time.Time
	if history != nil && !history.NewestDate.IsZero() {
		since = history.NewestDate
	} else {
		// Fetch from the beginning of Strava time (2008).
		since = time.Date(2008, 1, 1, 0, 0, 0, 0, time.UTC)
		history = &cache.HistoryData{}
	}

	newRuns, err := client.GetAllRunsSince(since)
	if err != nil {
		return fmt.Errorf("could not fetch runs for leaderboard: %w", err)
	}

	updated := cache.AppendHistory(history, newRuns)
	return cache.SaveHistory(updated)
}
