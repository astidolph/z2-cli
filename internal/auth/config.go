package auth

import (
	"github.com/z2-cli/internal/model"
	"github.com/z2-cli/internal/storage"
)

// Config is an alias for model.Config, preserving backward compatibility.
type Config = model.Config

func SaveConfig(config *Config) error {
	return storage.Get().SaveConfig(config)
}

func LoadConfig() (*Config, error) {
	return storage.Get().LoadConfig()
}
