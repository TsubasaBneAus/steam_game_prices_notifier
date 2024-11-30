package config

import (
	"context"
	"log/slog"

	"github.com/caarlos0/env/v11"
)

// A struct to store the configuration for Discord
type DiscordConfig struct {
	DiscordWebhookURL string `env:"DISCORD_WEBHOOK_URL,notEmpty"`
}

// Generate configuration for the Discord
func NewDiscordConfig(ctx context.Context) (*DiscordConfig, error) {
	cfg := &DiscordConfig{}
	if err := env.Parse(cfg); err != nil {
		slog.ErrorContext(
			ctx,
			"failed to load configuration for Discord",
			slog.Any("error", err),
		)

		return nil, err
	}

	return cfg, nil
}
