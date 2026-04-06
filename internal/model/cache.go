package model

import (
	"time"

	"github.com/z2-cli/internal/strava"
)

const DefaultTTL = 15 * time.Minute

type CachedData struct {
	FetchedAt  time.Time         `json:"fetched_at"`
	SinceUnix  int64             `json:"since_unix"`
	Activities []strava.Activity `json:"activities"`
}

// IsFresh returns true if the cached data covers the requested time range
// and was fetched within the TTL.
func (c *CachedData) IsFresh(since time.Time) bool {
	if time.Since(c.FetchedAt) > DefaultTTL {
		return false
	}
	return c.SinceUnix <= since.Unix()
}

type HistoryData struct {
	Activities []strava.Activity `json:"activities"`
	NewestDate time.Time         `json:"newest_date"`
}
