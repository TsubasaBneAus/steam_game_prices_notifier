//go:build wireinject

package main

import (
	"context"

	"github.com/TsubasaBneAus/steam_game_price_notifier/app/external/discord"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/external/httpclient"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/external/notion"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/external/steam"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/interactor"
	"github.com/TsubasaBneAus/steam_game_price_notifier/config"
	"github.com/google/wire"
)

// A wire set for the main package
var Set = wire.NewSet(
	NewApp,
	config.Set,
	httpclient.Set,
	steam.Set,
	notion.Set,
	discord.Set,
	interactor.Set,
)

// Initialize the application
func InitializeApp(ctx context.Context) (*app, error) {
	wire.Build(Set)
	return &app{}, nil
}
