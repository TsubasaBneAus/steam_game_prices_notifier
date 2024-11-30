package config

import (
	"context"
	"testing"
)

func TestNewDiscordConfig(t *testing.T) {
	t.Run("Positive case:  Successfully load configuration for Discord", func(t *testing.T) {
		// Set environment variables
		t.Setenv("DISCORD_WEBHOOK_URL", "dummy_discord_webhook_url")

		// Execute the function to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if _, err := NewDiscordConfig(ctx); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Negative case: Environment variables are missing or empty", func(t *testing.T) {
		// Set environment variables
		t.Setenv("DISCORD_WEBHOOK_URL", "")

		// Execute the function to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if _, err := NewDiscordConfig(ctx); err == nil {
			t.Errorf("\ngot: %v\nwant: an error generated in discord.go", nil)
		}
	})
}
