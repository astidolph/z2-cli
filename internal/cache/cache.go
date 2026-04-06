package cache

import (
	"github.com/z2-cli/internal/model"
	"github.com/z2-cli/internal/storage"
)

// CachedData is an alias for model.CachedData, preserving backward compatibility.
type CachedData = model.CachedData

// Load reads the cached data. Returns nil if missing or unparseable.
func Load() *CachedData {
	return storage.Get().LoadCache()
}

// Save writes cached data.
func Save(cached *CachedData) error {
	return storage.Get().SaveCache(cached)
}

// Invalidate deletes the cache.
func Invalidate() error {
	return storage.Get().InvalidateCache()
}
