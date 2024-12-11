package notion

import (
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/service"
	"github.com/google/wire"
)

// A wire set for the notion package
var Set = wire.NewSet(
	NewNotionWishlistGetter,
	NewNotionWishlistItemCreator,
	NewNotionWishlistItemUpdater,
	wire.Bind(new(service.NotionWishlistGetter), new(*notionWishlistGetter)),
	wire.Bind(new(service.NotionWishlistItemCreator), new(*notionWishlistItemCreator)),
	wire.Bind(new(service.NotionWishlistItemUpdater), new(*notionWishlistItemUpdater)),
)
