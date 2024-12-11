package discord

import (
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/service"
	"github.com/google/wire"
)

// A wire set for the discord package
var Set = wire.NewSet(
	NewVideoGamePricesOnDiscordNotifier,
	NewErrorOnDiscordNotifier,
	wire.Bind(new(service.VideoGamePricesOnDiscordNotifier), new(*videoGamePricesOnDiscordNotifier)),
	wire.Bind(new(service.ErrorOnDiscordNotifier), new(*errorOnDiscordNotifier)),
)
