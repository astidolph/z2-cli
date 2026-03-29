package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Zone2HR      int    `json:"zone2_hr,omitempty"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	dir := filepath.Join(home, ".z2-cli")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}
	return filepath.Join(dir, "config.json"), nil
}

func SaveConfig(config *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("could not save config: %w", err)
	}
	return nil
}

func LoadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not configured — run 'z2-cli auth' to set up your Strava API credentials")
		}
		return nil, fmt.Errorf("could not read config: %w", err)
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("could not parse config: %w", err)
	}
	return &config, nil
}
