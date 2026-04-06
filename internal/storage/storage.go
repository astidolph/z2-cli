package storage

import "github.com/z2-cli/internal/model"

// Store abstracts persistence for config, tokens, cache, and history.
// Implementations: FileStore (local JSON files) and PGStore (PostgreSQL).
type Store interface {
	LoadConfig() (*model.Config, error)
	SaveConfig(config *model.Config) error

	LoadToken() (*model.Token, error)
	SaveToken(token *model.Token) error

	LoadCache() *model.CachedData
	SaveCache(data *model.CachedData) error
	InvalidateCache() error

	LoadHistory() *model.HistoryData
	SaveHistory(data *model.HistoryData) error

	LoadSessionKey() ([]byte, error)
	SaveSessionKey(key []byte) error

	Close()
}
