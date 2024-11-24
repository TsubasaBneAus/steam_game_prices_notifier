package config

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type Envs struct {
	NotionApiKey      string
	NotionDatabaseID  string
	DiscordWebhookURL string
}

// Load environment variables
func LoadEnvs(ctx context.Context) (*Envs, error) {
	missingEnvs := make([]string, 0, 2)
	notionApiKey := os.Getenv("NOTION_API_KEY")
	if notionApiKey == "" {
		missingEnvs = append(missingEnvs, "NOTION_API_KEY")
	}
	notionDatabaseID := os.Getenv("NOTION_DATABASE_ID")
	if notionDatabaseID == "" {
		missingEnvs = append(missingEnvs, "NOTION_DATABASE_ID")
	}
	discordWebhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if discordWebhookURL == "" {
		missingEnvs = append(missingEnvs, "DISCORD_WEBHOOK_URL")
	}

	if len(missingEnvs) > 0 {
		return nil, fmt.Errorf("missing environment variables: %s", strings.Join(missingEnvs, ", "))
	}

	return &Envs{
		NotionApiKey:      notionApiKey,
		NotionDatabaseID:  notionDatabaseID,
		DiscordWebhookURL: discordWebhookURL,
	}, nil
}
