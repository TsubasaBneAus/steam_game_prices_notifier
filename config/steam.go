package config

import (
	"context"
	"log/slog"

	"github.com/caarlos0/env/v11"
)

// A struct to store the configuration for Steamworks API
type SteamConfig struct {
	SteamUserID string `env:"STEAM_USER_ID,notEmpty"`
}

// Generate configuration for the unofficial Steam API
func NewSteamConfig(ctx context.Context) (*SteamConfig, error) {
	cfg := &SteamConfig{}
	if err := env.Parse(cfg); err != nil {
		slog.ErrorContext(
			ctx,
			"failed to load configuration for the Notion API",
			slog.Any("error", err),
		)

		return nil, err
	}

	return cfg, nil
}
