package discord

import (
	"context"

	"github.com/TsubasaBneAus/steam_game_price_notifier/app/service"
)

type videoGamePricesOnDiscordNotifier struct {
	httpClient service.HTTPClient
}

var _ service.VideoGamePricesOnDiscordNotifier = (*videoGamePricesOnDiscordNotifier)(nil)

// Generate a new video game prices on Discord notifier
func NewVideoGamePricesOnDiscordNotifier(httpClient service.HTTPClient) *videoGamePricesOnDiscordNotifier {
	return &videoGamePricesOnDiscordNotifier{
		httpClient: httpClient,
	}
}

// Notify video game prices on Discord
func (n *videoGamePricesOnDiscordNotifier) NotifyVideoGamePricesOnDiscord(
	ctx context.Context,
	input *service.NotifyVideoGamePricesOnDiscordInput,
) (*service.NotifyVideoGamePricesOnDiscordOutput, error) {
	return &service.NotifyVideoGamePricesOnDiscordOutput{}, nil
}
