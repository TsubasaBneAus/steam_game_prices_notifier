package config

import (
	"context"
	"log/slog"

	"github.com/caarlos0/env/v11"
)

// A struct to store the configuration for Notion API
type NotionConfig struct {
	NotionAPIKey     string `env:"NOTION_API_KEY,notEmpty"`
	NotionDatabaseID string `env:"NOTION_DATABASE_ID,notEmpty"`
}

// Generate configuration for Notion API
func NewNotionConfig(ctx context.Context) (*NotionConfig, error) {
	cfg := &NotionConfig{}
	if err := env.Parse(cfg); err != nil {
		slog.ErrorContext(
			ctx,
			"failed to load configuration for Notion API",
			slog.Any("error", err),
		)

		return nil, err
	}

	return cfg, nil
}
