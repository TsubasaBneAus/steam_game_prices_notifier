package interactor

import (
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/usecase"
	"github.com/google/wire"
)

// A wire set for the interactor package
var Set = wire.NewSet(
	NewGamePricesNotifier,
	NewErrorOnDiscordNotifier,
	wire.Bind(new(usecase.VideoGamePricesNotifier), new(*videoGamePricesNotifier)),
	wire.Bind(new(usecase.ErrorNotifier), new(*errorOnDiscordNotifier)),
)
