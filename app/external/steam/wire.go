package steam

import (
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/service"
	"github.com/google/wire"
)

// A wire set for the steam package
var Set = wire.NewSet(
	NewSteamWishlistGetter,
	NewSteamVideoGameDetailsGetter,
	wire.Bind(new(service.SteamWishlistGetter), new(*steamWishlistGetter)),
	wire.Bind(new(service.SteamVideoGameDetailsGetter), new(*steamVideoGameDetailsGetter)),
)
