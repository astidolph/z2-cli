package cache

import (
	"sort"

	"github.com/z2-cli/internal/model"
	"github.com/z2-cli/internal/storage"
	"github.com/z2-cli/internal/strava"
)

// HistoryData is an alias for model.HistoryData, preserving backward compatibility.
type HistoryData = model.HistoryData

// LoadHistory reads the full history cache.
// Returns nil if not found or cannot be parsed.
func LoadHistory() *HistoryData {
	return storage.Get().LoadHistory()
}

// SaveHistory writes the full history cache.
func SaveHistory(history *HistoryData) error {
	return storage.Get().SaveHistory(history)
}

// AppendHistory merges new activities into the existing history, deduplicating
// by activity ID, and updates NewestDate.
func AppendHistory(existing *HistoryData, newRuns []strava.Activity) *HistoryData {
	seen := make(map[int64]bool, len(existing.Activities))
	for _, a := range existing.Activities {
		seen[a.ID] = true
	}

	merged := make([]strava.Activity, len(existing.Activities))
	copy(merged, existing.Activities)

	for _, a := range newRuns {
		if !seen[a.ID] {
			merged = append(merged, a)
			seen[a.ID] = true
		}
	}

	// Find the newest date across all activities.
	newest := existing.NewestDate
	for _, a := range merged {
		if t, err := a.StartTime(); err == nil && t.After(newest) {
			newest = t
		}
	}

	// Sort by start date descending for consistent ordering.
	sort.Slice(merged, func(i, j int) bool {
		ti, _ := merged[i].StartTime()
		tj, _ := merged[j].StartTime()
		return ti.After(tj)
	})

	return &HistoryData{
		Activities: merged,
		NewestDate: newest,
	}
}
