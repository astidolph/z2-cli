package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/z2-cli/internal/strava"
)

type HistoryData struct {
	Activities []strava.Activity `json:"activities"`
	// NewestDate is the start time of the most recent cached activity,
	// used as the "after" parameter for incremental syncs.
	NewestDate time.Time `json:"newest_date"`
}

func historyPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	dir := filepath.Join(home, ".z2-cli")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}
	return filepath.Join(dir, "history.json"), nil
}

// LoadHistory reads the full history cache from disk.
// Returns nil if the file does not exist or cannot be parsed.
func LoadHistory() *HistoryData {
	path, err := historyPath()
	if err != nil {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var history HistoryData
	if err := json.Unmarshal(data, &history); err != nil {
		return nil
	}
	return &history
}

// SaveHistory writes the full history cache to disk.
func SaveHistory(history *HistoryData) error {
	path, err := historyPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal history: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("could not save history: %w", err)
	}
	return nil
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
