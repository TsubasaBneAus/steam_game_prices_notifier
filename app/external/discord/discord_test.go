package discord

import (
	"context"
	"errors"
	"net/http"
	"testing"

	httpclient "github.com/TsubasaBneAus/steam_game_price_notifier/app/external/httpclient/mock"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/model"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/service"
	"github.com/TsubasaBneAus/steam_game_price_notifier/config"
	"go.uber.org/mock/gomock"
)

func TestNotifyVideoGamePricesOnDiscord(t *testing.T) {
	t.Parallel()

	t.Run("Positive case: Successfully notify video game prices on Discord", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(
				&http.Response{
					StatusCode: http.StatusNoContent,
					Body:       http.NoBody,
				}, nil,
			)

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.DiscordConfig{
			DiscordWebhookID:    "dummy_discord_webhook_id",
			DiscordWebhookToken: "dummy_discord_webhook_token",
		}
		n := NewVideoGamePricesOnDiscordNotifier(cfg, m)
		input := &service.NotifyVideoGamePricesOnDiscordInput{
			DiscordContents: map[model.SteamAppID]*model.DiscordContent{
				1: {
					Title:        "dummy_title",
					CurrentPrice: 1000,
					LowestPrice:  1500,
				},
			},
		}
		if _, err := n.NotifyVideoGamePricesOnDiscord(ctx, input); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Positive case: Successfully notify video game prices on Discord with empty input", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(
				&http.Response{
					StatusCode: http.StatusNoContent,
					Body:       http.NoBody,
				}, nil,
			)

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.DiscordConfig{
			DiscordWebhookID:    "dummy_discord_webhook_id",
			DiscordWebhookToken: "dummy_discord_webhook_token",
		}
		n := NewVideoGamePricesOnDiscordNotifier(cfg, m)
		input := &service.NotifyVideoGamePricesOnDiscordInput{
			DiscordContents: map[model.SteamAppID]*model.DiscordContent{},
		}
		if _, err := n.NotifyVideoGamePricesOnDiscord(ctx, input); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Negative case: Failed to send a Discord API request", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		wantErr := errors.New("unexpected error")
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(nil, wantErr)

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.DiscordConfig{
			DiscordWebhookID:    "dummy_discord_webhook_id",
			DiscordWebhookToken: "dummy_discord_webhook_token",
		}
		n := NewVideoGamePricesOnDiscordNotifier(cfg, m)
		input := &service.NotifyVideoGamePricesOnDiscordInput{
			DiscordContents: map[model.SteamAppID]*model.DiscordContent{
				1: {
					Title:        "dummy_title",
					CurrentPrice: 1000,
					LowestPrice:  1500,
				},
			},
		}
		if _, gotErr := n.NotifyVideoGamePricesOnDiscord(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Get a status code except 200", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(
				&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       http.NoBody,
				}, nil,
			)

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.DiscordConfig{
			DiscordWebhookID:    "dummy_discord_webhook_id",
			DiscordWebhookToken: "dummy_discord_webhook_token",
		}
		n := NewVideoGamePricesOnDiscordNotifier(cfg, m)
		input := &service.NotifyVideoGamePricesOnDiscordInput{
			DiscordContents: map[model.SteamAppID]*model.DiscordContent{
				1: {
					Title:        "dummy_title",
					CurrentPrice: 1000,
					LowestPrice:  1500,
				},
			},
		}
		wantErr := errUnexpectedStatusCode
		if _, gotErr := n.NotifyVideoGamePricesOnDiscord(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})
}

func TestNotifyErrorOnDiscord(t *testing.T) {
	t.Parallel()

	t.Run("Positive case: Successfully notify an error on Discord", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(
				&http.Response{
					StatusCode: http.StatusNoContent,
					Body:       http.NoBody,
				}, nil,
			)

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.DiscordConfig{
			DiscordWebhookID:    "dummy_discord_webhook_id",
			DiscordWebhookToken: "dummy_discord_webhook_token",
		}
		n := NewErrorOnDiscordNotifier(cfg, m)
		input := &service.NotifyErrorOnDiscordInput{
			GeneratedError: errors.New("dummy_error"),
		}
		if _, err := n.NotifyErrorOnDiscord(ctx, input); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Negative case: Failed to send a Discord API request", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		wantErr := errors.New("unexpected error")
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(nil, wantErr)

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.DiscordConfig{
			DiscordWebhookID:    "dummy_discord_webhook_id",
			DiscordWebhookToken: "dummy_discord_webhook_token",
		}
		n := NewErrorOnDiscordNotifier(cfg, m)
		input := &service.NotifyErrorOnDiscordInput{
			GeneratedError: errors.New("dummy_error"),
		}
		if _, gotErr := n.NotifyErrorOnDiscord(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Get a status code except 200", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(
				&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       http.NoBody,
				}, nil,
			)

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.DiscordConfig{
			DiscordWebhookID:    "dummy_discord_webhook_id",
			DiscordWebhookToken: "dummy_discord_webhook_token",
		}
		n := NewErrorOnDiscordNotifier(cfg, m)
		input := &service.NotifyErrorOnDiscordInput{
			GeneratedError: errors.New("dummy_error"),
		}
		wantErr := errUnexpectedStatusCode
		if _, gotErr := n.NotifyErrorOnDiscord(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})
}
