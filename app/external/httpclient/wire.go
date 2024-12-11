package httpclient

import (
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/service"
	"github.com/google/wire"
)

// A wire set for the httpclient package
var Set = wire.NewSet(
	NewHTTPClient,
	wire.Bind(new(service.HTTPClient), new(*httpClient)),
)
