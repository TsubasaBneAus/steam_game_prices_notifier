package config

import (
	"context"
	"testing"
)

func TestNewNotionConfig(t *testing.T) {
	t.Run("Positive case: Successfully load configuration for Notion API", func(t *testing.T) {
		// Set environment variables
		t.Setenv("NOTION_API_KEY", "dummy_notion_api_key")
		t.Setenv("NOTION_DATABASE_ID", "dummy_notion_database_id")

		// Execute the function to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if _, err := NewNotionConfig(ctx); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Negative case: Environment variables are missing or empty", func(t *testing.T) {
		// Set environment variables
		t.Setenv("NOTION_API_KEY", "")
		t.Setenv("NOTION_DATABASE_ID", "")

		// Execute the function to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if _, err := NewNotionConfig(ctx); err == nil {
			t.Errorf("\ngot: %v\nwant: an error generated in notion.go", nil)
		}
	})
}
