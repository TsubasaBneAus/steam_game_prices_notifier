package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/TsubasaBneAus/steam_game_price_notifier/app/model"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/service"
	"github.com/TsubasaBneAus/steam_game_price_notifier/config"
	"golang.org/x/time/rate"
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
//
// [FYI]
// The rate limiter is set to 5 requests per second and parallel processing is used
// ref. https://discord.com/developers/docs/topics/rate-limits
func (n *videoGamePricesOnDiscordNotifier) NotifyVideoGamePricesOnDiscord(
	ctx context.Context,
	input *service.NotifyVideoGamePricesOnDiscordInput,
) (*service.NotifyVideoGamePricesOnDiscordOutput, error) {
	limiter := rate.NewLimiter(5, 1)
	for _, v := range n.buildMessageBody(input.DiscordContents) {
		if err := limiter.Wait(ctx); err != nil {
			slog.ErrorContext(ctx, "failed to wait for the rate limiter", slog.Any("error", err))
			return nil, err
		}

		body := &model.DiscordMessageBody{
			Content: strings.Join(v, "\n"),
		}

		if err := n.notifyVideoGamePricesOnDiscord(ctx, body); err != nil {
			return nil, err
		}
	}

	return &service.NotifyVideoGamePricesOnDiscordOutput{}, nil
}

func (n *videoGamePricesOnDiscordNotifier) notifyVideoGamePricesOnDiscord(
	ctx context.Context,
	body *model.DiscordMessageBody,
) error {
	reqURL, err := url.JoinPath(discordAPIURL, "webhooks", n.cfg.DiscordWebhookID, n.cfg.DiscordWebhookToken)
	if err != nil {
		slog.ErrorContext(ctx, "failed to build a Discord API URL", slog.Any("error", err))
		return err
	}

	reqJSON, err := json.Marshal(body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal a Notion API request body", slog.Any("error", err))
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(reqJSON))
	if err != nil {
		slog.ErrorContext(ctx, "failed to create a Discord API request", slog.Any("error", err))
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := n.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send a Discord API request", slog.Any("error", err))
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		slog.ErrorContext(ctx, "failed to send a Discord API request", slog.Any("status_code", res.StatusCode))
		return errUnexpectedStatusCode
	}

	return nil
}

// Build a message body of a Discord message
func (n *videoGamePricesOnDiscordNotifier) buildMessageBody(
	discordContents map[model.SteamAppID]*model.DiscordContent,
) [][]string {
	contents := make(map[string]string, len(discordContents))
	for _, v := range discordContents {
		content := fmt.Sprintf(
			"- Title: **%s**  |  Current Price: **%d (JPY)**  |  Lowest Price: **%d (JPY)**",
			v.Title,
			v.CurrentPrice,
			v.LowestPrice,
		)
		contents[v.Title] = content
	}

	// Sort the contents by a video game title in ascending order
	//
	// [FYI]
	// The Discord message has a limitation of 2000 characters
	// Therefore, the contents are divided into multiple messages by 10 video games
	sortedContents := make([]string, 0, 10)
	contentsList := make([][]string, 0, len(contents))
	sortedContents = append(sortedContents, "## The recommended video games to buy now are as follows:")
	var count uint8
	for _, k := range slices.Sorted(maps.Keys(contents)) {
		if count == 10 {
			contentsList = append(contentsList, sortedContents)
			sortedContents = nil
			count = 0
		}

		sortedContents = append(sortedContents, contents[k])
		count++
	}
	contentsList = append(contentsList, sortedContents)

	return contentsList
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
		Content: fmt.Sprintf("## An error occurred:\n- %s", input.GeneratedError.Error()),
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
