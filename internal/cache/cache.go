package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/z2-cli/internal/strava"
)

const defaultTTL = 15 * time.Minute

type CachedData struct {
	FetchedAt  time.Time          `json:"fetched_at"`
	SinceUnix  int64              `json:"since_unix"`
	Activities []strava.Activity  `json:"activities"`
}

func cachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	dir := filepath.Join(home, ".z2-cli")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}
	return filepath.Join(dir, "cache.json"), nil
}

// Load reads the cached data from disk. Returns nil if the cache file
// does not exist or cannot be parsed.
func Load() *CachedData {
	path, err := cachePath()
	if err != nil {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var cached CachedData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil
	}
	return &cached
}

// Save writes cached data to disk.
func Save(cached *CachedData) error {
	path, err := cachePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cached, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal cache: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("could not save cache: %w", err)
	}
	return nil
}

// IsFresh returns true if the cached data covers the requested time range
// and was fetched within the TTL.
func (c *CachedData) IsFresh(since time.Time) bool {
	if time.Since(c.FetchedAt) > defaultTTL {
		return false
	}
	// Cache must cover at least as far back as requested
	return c.SinceUnix <= since.Unix()
}

// Invalidate deletes the cache file.
func Invalidate() error {
	path, err := cachePath()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
