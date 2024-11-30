package config

import (
	"context"
	"testing"
)

func TestNewSteamConfig(t *testing.T) {
	t.Run("Positive case: Successfully load configuration for an unofficial Steam API", func(t *testing.T) {
		// Set environment variables
		t.Setenv("STEAM_USER_ID", "dummy_steam_user_id")

		// Execute the function to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if _, err := NewSteamConfig(ctx); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Negative case: Environment variables are missing or empty", func(t *testing.T) {
		// Set environment variables
		t.Setenv("STEAM_USER_ID", "")

		// Execute the function to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if _, err := NewSteamConfig(ctx); err == nil {
			t.Errorf("\ngot: %v\nwant: an error generated in steam.go", nil)
		}
	})
}
