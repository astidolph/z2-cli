package storage

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/z2-cli/internal/model"
)

// FileStore persists data as JSON files in ~/.z2-cli/.
type FileStore struct {
	dir string
}

// NewFileStore creates a FileStore using ~/.z2-cli/ as the storage directory.
func NewFileStore() (*FileStore, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not find home directory: %w", err)
	}
	dir := filepath.Join(home, ".z2-cli")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("could not create config directory: %w", err)
	}
	return &FileStore{dir: dir}, nil
}

func (f *FileStore) path(name string) string {
	return filepath.Join(f.dir, name)
}

func (f *FileStore) Close() {}

// --- Config ---

func (f *FileStore) LoadConfig() (*model.Config, error) {
	data, err := os.ReadFile(f.path("config.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not configured — run 'z2-cli auth' to set up your Strava API credentials")
		}
		return nil, fmt.Errorf("could not read config: %w", err)
	}
	var config model.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("could not parse config: %w", err)
	}
	return &config, nil
}

func (f *FileStore) SaveConfig(config *model.Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}
	if err := os.WriteFile(f.path("config.json"), data, 0600); err != nil {
		return fmt.Errorf("could not save config: %w", err)
	}
	return nil
}

// --- Token ---

func (f *FileStore) LoadToken() (*model.Token, error) {
	data, err := os.ReadFile(f.path("token.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not authenticated — run 'z2-cli auth' to connect your Strava account")
		}
		return nil, fmt.Errorf("could not read token: %w", err)
	}
	var token model.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}
	return &token, nil
}

func (f *FileStore) SaveToken(token *model.Token) error {
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal token: %w", err)
	}
	if err := os.WriteFile(f.path("token.json"), data, 0600); err != nil {
		return fmt.Errorf("could not save token: %w", err)
	}
	return nil
}

// --- Cache ---

func (f *FileStore) LoadCache() *model.CachedData {
	data, err := os.ReadFile(f.path("cache.json"))
	if err != nil {
		return nil
	}
	var cached model.CachedData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil
	}
	return &cached
}

func (f *FileStore) SaveCache(cached *model.CachedData) error {
	data, err := json.MarshalIndent(cached, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal cache: %w", err)
	}
	if err := os.WriteFile(f.path("cache.json"), data, 0600); err != nil {
		return fmt.Errorf("could not save cache: %w", err)
	}
	return nil
}

func (f *FileStore) InvalidateCache() error {
	err := os.Remove(f.path("cache.json"))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// --- History ---

func (f *FileStore) LoadHistory() *model.HistoryData {
	data, err := os.ReadFile(f.path("history.json"))
	if err != nil {
		return nil
	}
	var history model.HistoryData
	if err := json.Unmarshal(data, &history); err != nil {
		return nil
	}
	return &history
}

func (f *FileStore) SaveHistory(history *model.HistoryData) error {
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal history: %w", err)
	}
	if err := os.WriteFile(f.path("history.json"), data, 0600); err != nil {
		return fmt.Errorf("could not save history: %w", err)
	}
	return nil
}

// --- Session Key ---

func (f *FileStore) LoadSessionKey() ([]byte, error) {
	data, err := os.ReadFile(f.path("session_key"))
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strings.TrimSpace(string(data)))
}

func (f *FileStore) SaveSessionKey(key []byte) error {
	return os.WriteFile(f.path("session_key"), []byte(hex.EncodeToString(key)), 0600)
}
