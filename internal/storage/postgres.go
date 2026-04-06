package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/z2-cli/internal/model"
)

// PGStore persists data in a PostgreSQL database using a simple key-value table.
type PGStore struct {
	pool *pgxpool.Pool
}

// NewPGStore connects to PostgreSQL and ensures the kv table exists.
func NewPGStore(ctx context.Context, databaseURL string) (*PGStore, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS kv (
			key        TEXT PRIMARY KEY,
			value      JSONB NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("could not create kv table: %w", err)
	}

	return &PGStore{pool: pool}, nil
}

func (p *PGStore) Close() {
	p.pool.Close()
}

func (p *PGStore) load(key string, dest any) error {
	var raw []byte
	err := p.pool.QueryRow(context.Background(),
		`SELECT value FROM kv WHERE key = $1`, key).Scan(&raw)
	if err == pgx.ErrNoRows {
		return err
	}
	if err != nil {
		return fmt.Errorf("could not read %s: %w", key, err)
	}
	return json.Unmarshal(raw, dest)
}

func (p *PGStore) save(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("could not marshal %s: %w", key, err)
	}
	_, err = p.pool.Exec(context.Background(),
		`INSERT INTO kv (key, value, updated_at) VALUES ($1, $2, now())
		 ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = now()`,
		key, data)
	if err != nil {
		return fmt.Errorf("could not save %s: %w", key, err)
	}
	return nil
}

// --- Config ---

func (p *PGStore) LoadConfig() (*model.Config, error) {
	var config model.Config
	err := p.load("config", &config)
	if err == pgx.ErrNoRows {
		// In PG mode, credentials come from env vars. Zone2HR may not be set yet.
		config = model.Config{}
	} else if err != nil {
		return nil, err
	}

	// Overlay credentials from environment variables.
	if id := os.Getenv("STRAVA_CLIENT_ID"); id != "" {
		config.ClientID = id
	}
	if secret := os.Getenv("STRAVA_CLIENT_SECRET"); secret != "" {
		config.ClientSecret = secret
	}

	if config.ClientID == "" || config.ClientSecret == "" {
		return nil, fmt.Errorf("not configured — set STRAVA_CLIENT_ID and STRAVA_CLIENT_SECRET environment variables")
	}

	return &config, nil
}

func (p *PGStore) SaveConfig(config *model.Config) error {
	// Only persist Zone2HR to the database; credentials stay in env vars.
	dbConfig := model.Config{Zone2HR: config.Zone2HR}
	return p.save("config", &dbConfig)
}

// --- Token ---

func (p *PGStore) LoadToken() (*model.Token, error) {
	var token model.Token
	err := p.load("token", &token)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("not authenticated — use the web UI to connect your Strava account")
	}
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (p *PGStore) SaveToken(token *model.Token) error {
	return p.save("token", token)
}

// --- Cache ---

func (p *PGStore) LoadCache() *model.CachedData {
	var cached model.CachedData
	if err := p.load("cache", &cached); err != nil {
		return nil
	}
	return &cached
}

func (p *PGStore) SaveCache(cached *model.CachedData) error {
	return p.save("cache", cached)
}

func (p *PGStore) InvalidateCache() error {
	_, err := p.pool.Exec(context.Background(),
		`DELETE FROM kv WHERE key = 'cache'`)
	return err
}

// --- History ---

func (p *PGStore) LoadHistory() *model.HistoryData {
	var history model.HistoryData
	if err := p.load("history", &history); err != nil {
		return nil
	}
	return &history
}

func (p *PGStore) SaveHistory(history *model.HistoryData) error {
	return p.save("history", history)
}
