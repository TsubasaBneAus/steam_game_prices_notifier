package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/TsubasaBneAus/steam_game_price_notifier/app/model"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/service"
	"github.com/TsubasaBneAus/steam_game_price_notifier/config"
)

const discordAPIURL string = "https://discord.com/api"

type videoGamePricesOnDiscordNotifier struct {
	cfg        *config.DiscordConfig
	httpClient service.HTTPClient
}

var _ service.VideoGamePricesOnDiscordNotifier = (*videoGamePricesOnDiscordNotifier)(nil)

// Generate a new video game prices on Discord notifier
func NewVideoGamePricesOnDiscordNotifier(
	cfg *config.DiscordConfig,
	httpClient service.HTTPClient,
) *videoGamePricesOnDiscordNotifier {
	return &videoGamePricesOnDiscordNotifier{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

// Notify video game prices on Discord
func (n *videoGamePricesOnDiscordNotifier) NotifyVideoGamePricesOnDiscord(
	ctx context.Context,
	input *service.NotifyVideoGamePricesOnDiscordInput,
) (*service.NotifyVideoGamePricesOnDiscordOutput, error) {
	reqURL, err := url.JoinPath(discordAPIURL, "webhooks", n.cfg.DiscordWebhookID, n.cfg.DiscordWebhookToken)
	if err != nil {
		slog.ErrorContext(ctx, "failed to build a Discord API URL", slog.Any("error", err))
		return nil, err
	}

	// Build a request body of a Discord message
	contents := make([]string, 0, len(input.DiscordContents))
	contents = append(contents, "## The recommended video games to buy now are as follows:")
	for _, v := range input.DiscordContents {
		content := fmt.Sprintf(
			"- Title: **%s**  |  Current Price: **%d (JPY)**  |  Lowest Price: **%d (JPY)**",
			v.Title,
			v.CurrentPrice,
			v.LowestPrice,
		)
		contents = append(contents, content)
	}
	body := &model.DiscordMessageBody{
		Content: strings.Join(contents, "\n"),
	}
	reqJSON, err := json.Marshal(body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal a Notion API request body", slog.Any("error", err))
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(reqJSON))
	if err != nil {
		slog.ErrorContext(ctx, "failed to create a Discord API request", slog.Any("error", err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := n.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send a Discord API request", slog.Any("error", err))
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		slog.ErrorContext(ctx, "failed to send a Discord API request", slog.Any("status_code", res.StatusCode))
		return nil, errUnexpectedStatusCode
	}

	return &service.NotifyVideoGamePricesOnDiscordOutput{}, nil
}

type errorOnDiscordNotifier struct {
	cfg        *config.DiscordConfig
	httpClient service.HTTPClient
}

var _ service.ErrorOnDiscordNotifier = (*errorOnDiscordNotifier)(nil)

// Generate a new error on Discord notifier
func NewErrorOnDiscordNotifier(
	cfg *config.DiscordConfig,
	httpClient service.HTTPClient,
) *errorOnDiscordNotifier {
	return &errorOnDiscordNotifier{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

// Notify an error on Discord
func (n *errorOnDiscordNotifier) NotifyErrorOnDiscord(
	ctx context.Context,
	input *service.NotifyErrorOnDiscordInput,
) (*service.NotifyErrorOnDiscordOutput, error) {
	reqURL, err := url.JoinPath(discordAPIURL, "webhooks", n.cfg.DiscordWebhookID, n.cfg.DiscordWebhookToken)
	if err != nil {
		slog.ErrorContext(ctx, "failed to build a Discord API URL", slog.Any("error", err))
		return nil, err
	}

	// Build a request body of a Discord message
	body := &model.DiscordMessageBody{
		Content: fmt.Sprintf("## An error occurred:\n%v", input.GeneratedError),
	}
	reqJSON, err := json.Marshal(body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal a Notion API request body", slog.Any("error", err))
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(reqJSON))
	if err != nil {
		slog.ErrorContext(ctx, "failed to create a Discord API request", slog.Any("error", err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := n.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send a Discord API request", slog.Any("error", err))
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		slog.ErrorContext(ctx, "failed to send a Discord API request", slog.Any("status_code", res.StatusCode))
		return nil, errUnexpectedStatusCode
	}

	return &service.NotifyErrorOnDiscordOutput{}, nil
}
