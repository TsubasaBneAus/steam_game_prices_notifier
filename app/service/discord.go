package service

import (
	"context"

	"github.com/TsubasaBneAus/steam_game_price_notifier/app/model"
)

//go:generate mockgen -source=./discord.go -destination=../external/discord/mock/discord.go -package=mock -typed

type (
	// An input to notify video game prices on Discord
	NotifyVideoGamePricesOnDiscordInput struct {
		DiscordContents map[model.SteamAppID]*model.DiscordContent
	}

	// An output to notify video game prices on Discord
	NotifyVideoGamePricesOnDiscordOutput struct{}

	// An interface to notify video game prices on Discord
	VideoGamePricesOnDiscordNotifier interface {
		NotifyVideoGamePricesOnDiscord(
			ctx context.Context,
			input *NotifyVideoGamePricesOnDiscordInput,
		) (*NotifyVideoGamePricesOnDiscordOutput, error)
	}
)

type (
	// An input to notify an error on Discord
	NotifyErrorOnDiscordInput struct {
		GeneratedError error
	}

	// An output to notify an error on Discord
	NotifyErrorOnDiscordOutput struct{}

	// An interface to notify an error on Discord
	ErrorOnDiscordNotifier interface {
		NotifyErrorOnDiscord(
			ctx context.Context,
			input *NotifyErrorOnDiscordInput,
		) (*NotifyErrorOnDiscordOutput, error)
	}
)
