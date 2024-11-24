package config

import (
	"context"
	"testing"
)

func TestLoadEnvs(t *testing.T) {
	t.Run("Positive case: Load environment variables successfully", func(t *testing.T) {
		// Set environment variables
		t.Setenv("NOTION_API_KEY", "dummy-notion-api-key")
		t.Setenv("NOTION_DATABASE_ID", "dummy-notion-database-id")
		t.Setenv("DISCORD_WEBHOOK_URL", "dummy-discord-webhook-url")

		// Execute the function to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if _, err := LoadEnvs(ctx); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Negative case: Environment variables are missing or empty", func(t *testing.T) {
		// Set environment variables
		t.Setenv("NOTION_API_KEY", "")
		t.Setenv("NOTION_DATABASE_ID", "")
		t.Setenv("DISCORD_WEBHOOK_URL", "")

		// Execute the function to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if _, err := LoadEnvs(ctx); err == nil {
			t.Errorf("\ngot: %v\nwant: an error generated in config.go", nil)
		}
	})
}
