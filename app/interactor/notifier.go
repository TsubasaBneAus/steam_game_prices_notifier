package interactor

import (
	"context"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/TsubasaBneAus/steam_game_price_notifier/app/model"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/service"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/usecase"
	"github.com/TsubasaBneAus/steam_game_price_notifier/config"
	multierror "github.com/hashicorp/go-multierror"
	"golang.org/x/time/rate"
)

type videoGamePricesNotifier struct {
	cfg           *config.NotionConfig
	sWGetter      service.SteamWishlistGetter
	sVGDGetter    service.SteamVideoGameDetailsGetter
	nWGetter      service.NotionWishlistGetter
	nWICreator    service.NotionWishlistItemCreator
	nWIUpdater    service.NotionWishlistItemUpdater
	nWIDeleter    service.NotionWishlistItemDeleter
	vGPODNotifier service.VideoGamePricesOnDiscordNotifier
}

var _ usecase.VideoGamePricesNotifier = (*videoGamePricesNotifier)(nil)

// Generate a new videoGamePricesNotifier
func NewGamePricesNotifier(
	cfg *config.NotionConfig,
	sWGetter service.SteamWishlistGetter,
	sVGDGetter service.SteamVideoGameDetailsGetter,
	nWGetter service.NotionWishlistGetter,
	nWICreator service.NotionWishlistItemCreator,
	nWIUpdater service.NotionWishlistItemUpdater,
	nWIDeleter service.NotionWishlistItemDeleter,
	vGPODNotifier service.VideoGamePricesOnDiscordNotifier,
) *videoGamePricesNotifier {
	return &videoGamePricesNotifier{
		cfg:           cfg,
		sWGetter:      sWGetter,
		sVGDGetter:    sVGDGetter,
		nWGetter:      nWGetter,
		nWICreator:    nWICreator,
		nWIUpdater:    nWIUpdater,
		nWIDeleter:    nWIDeleter,
		vGPODNotifier: vGPODNotifier,
	}
}

// Notify video game prices on Discord
func (n *videoGamePricesNotifier) NotifyVideoGamePrices(
	ctx context.Context,
	input *usecase.NotifyVideoGamePricesInput,
) (*usecase.NotifyVideoGamePricesOutput, error) {
	// Get a list of video game details on the Steam Store
	vGDList, err := n.getVideoGameDetailsList(ctx)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed to get a list of video game details on the Steam Store",
			slog.Any("error", err),
		)
		return nil, err
	}

	// Get a wishlist from the Notion DB
	nWishlist, err := n.nWGetter.GetNotionWishlist(ctx, &service.GetNotionWishlistInput{})
	if err != nil {
		slog.ErrorContext(ctx, "failed to get a Notion DB wishlist", slog.Any("error", err))
		return nil, err
	}

	// Create or update a wishlist on the Notion DB based on the Steam Store wishlist
	discordContents, err := n.createOrUpdateNotionWishlist(ctx, vGDList, nWishlist.WishlistItems)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create or update a wishlist on the Notion DB", slog.Any("error", err))
		return nil, err
	}

	// Delete a wishlist on the Notion DB
	if err := n.deleteNotionWishlistItems(ctx, vGDList, nWishlist.WishlistItems); err != nil {
		slog.ErrorContext(ctx, "failed to delete a wishlist on the Notion DB", slog.Any("error", err))
		return nil, err
	}

	// Terminate processing if there are no video game prices to notify
	if len(discordContents) == 0 {
		return &usecase.NotifyVideoGamePricesOutput{}, nil
	}

	// Notify video game prices on Discord
	vGPODInput := &service.NotifyVideoGamePricesOnDiscordInput{
		DiscordContents: discordContents,
	}
	if _, err := n.vGPODNotifier.NotifyVideoGamePricesOnDiscord(ctx, vGPODInput); err != nil {
		slog.ErrorContext(ctx, "failed to notify video game prices on Discord", slog.Any("error", err))
		return nil, err
	}

	return &usecase.NotifyVideoGamePricesOutput{}, nil
}

// Get a list of video game details on the Steam Store
func (n *videoGamePricesNotifier) getVideoGameDetailsList(
	ctx context.Context,
) (map[model.SteamAppID]*model.SteamStoreVideoGameDetails, error) {
	// Get a wishlist from the Steam Store
	steamWishlist, err := n.sWGetter.GetSteamWishlist(ctx, &service.GetSteamWishlistInput{})
	if err != nil {
		slog.ErrorContext(ctx, "failed to get a Steam Store wishlist", slog.Any("error", err))
		return nil, err
	}

	// Get a list of video game details on the Steam Store
	//
	// [FYI]
	// The rate limiter is set to 5 requests per second and parallel processing is used,
	// but there is no recommended rate limit due to the unofficial Steam Store API
	videoGameDetailsList := make(
		map[model.SteamAppID]*model.SteamStoreVideoGameDetails,
		len(steamWishlist.Wishlist.Response.Items),
	)
	var mu sync.Mutex
	limiter := rate.NewLimiter(5, 1)
	meg := &multierror.Group{}
	for _, item := range steamWishlist.Wishlist.Response.Items {
		if err := limiter.Wait(ctx); err != nil {
			slog.ErrorContext(ctx, "failed to wait the rate limiter", slog.Any("error", err))
			return nil, err
		}

		meg.Go(func() error {
			mu.Lock()
			defer mu.Unlock()

			// Get a video game details on the Steam Store
			input := &service.GetSteamVideoGameDetailsInput{
				AppID: model.SteamAppID(item.AppID),
			}
			videoGameDetails, err := n.sVGDGetter.GetSteamVideoGameDetails(ctx, input)
			if err != nil {
				slog.ErrorContext(
					ctx,
					"failed to get a video game details on the Steam Store",
					slog.Any("error", err),
				)
				return err
			}

			videoGameDetailsList[model.SteamAppID(item.AppID)] = videoGameDetails.VideoGameDetails

			return nil
		})
	}

	if err := meg.Wait(); err != nil {
		slog.ErrorContext(
			ctx,
			"failed to get a list of video game details on the Steam Store",
			slog.Any("error", err),
		)
		return nil, err
	}

	return videoGameDetailsList, nil
}

// Create or update a wishlist on the Notion DB based on the Steam Store wishlist
func (n *videoGamePricesNotifier) createOrUpdateNotionWishlist(
	ctx context.Context,
	vGDList map[model.SteamAppID]*model.SteamStoreVideoGameDetails,
	nWishList []*model.NotionWishlistItem,
) (map[model.SteamAppID]*model.DiscordContent, error) {
	// Convert a Notion wishlist to a map
	convertedNWishList := make(map[model.SteamAppID]*model.NotionWishlistItem, len(nWishList))
	for _, v := range nWishList {
		appID, err := strconv.Atoi(v.Properties.NotionAppID.Title[0].NotionText.NotionContent)
		if err != nil {
			slog.ErrorContext(ctx, "failed to convert the app ID to int", slog.Any("error", err))
			return nil, err
		}

		convertedNWishList[model.SteamAppID(appID)] = v
	}

	// Separate the video game details list into two lists: one to create and one to update
	listToCreate := make(map[model.SteamAppID]*model.SteamStoreVideoGameDetails, len(vGDList))
	listToUpdate := make(map[model.SteamAppID]*model.SteamStoreVideoGameDetails, len(vGDList))
	for i, v := range vGDList {
		if _, ok := convertedNWishList[i]; ok {
			listToUpdate[i] = v
		} else {
			listToCreate[i] = v
		}
	}

	// Create wishlist items on the Notion DB
	if err := n.createNotionWishlistItems(ctx, listToCreate); err != nil {
		slog.ErrorContext(ctx, "failed to create a wishlist item on the Notion DB", slog.Any("error", err))
		return nil, err
	}

	// Update wishlist items on the Notion DB
	discordContents, err := n.updateNotionWishlistItems(ctx, convertedNWishList, listToUpdate)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update a wishlist item on the Notion DB", slog.Any("error", err))
		return nil, err
	}

	return discordContents, nil
}

// Create a wishlist on the Notion DB
//
// [FYI]
// The rate limiter is set to 3 requests per second and parallel processing is used
func (n *videoGamePricesNotifier) createNotionWishlistItems(
	ctx context.Context,
	listToCreate map[model.SteamAppID]*model.SteamStoreVideoGameDetails,
) error {
	limiter := rate.NewLimiter(3, 1)
	meg := &multierror.Group{}
	for i, v := range listToCreate {
		if err := limiter.Wait(ctx); err != nil {
			slog.ErrorContext(ctx, "failed to wait the rate limiter", slog.Any("error", err))
			return err
		}

		meg.Go(func() error {
			// Convert the current price of a video game to uint64
			currentPrice, err := n.convertCurrentPrice(ctx, v.CurrentPrice)
			if err != nil {
				return err
			}

			input := &service.CreateNotionWishlistItemInput{
				WishlistItem: &model.NotionWishlistItem{
					Parent: &model.NotionParent{
						DatabaseID: model.NotionDatabaseID(n.cfg.NotionDatabaseID),
					},
					Properties: &model.NotionProperties{
						NotionAppID: &model.NotionAppID{
							Title: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: strconv.Itoa(int(i)),
									},
								},
							},
						},
						NotionName: &model.NotionName{
							RichText: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: v.Title,
									},
								},
							},
						},
						CurrentPrice: &model.NotionPrice{
							Number: currentPrice,
						},
						LowestPrice: &model.NotionPrice{
							Number: nil,
						},
						NotionReleaseDate: &model.NotionReleaseDate{
							NotionDate: n.convertReleaseDate(ctx, v.ReleaseDate),
						},
					},
				},
			}
			if _, err := n.nWICreator.CreateNotionWishlistItem(ctx, input); err != nil {
				slog.ErrorContext(ctx, "failed to create a wishlist item on the Notion DB", slog.Any("error", err))
				return err
			}

			return nil
		})
	}

	if err := meg.Wait(); err != nil {
		slog.ErrorContext(ctx, "failed to create a wishlist item on the Notion DB", slog.Any("error", err))
		return err
	}

	return nil
}

// Update a wishlist on the Notion DB
//
// [FYI]
// The rate limiter is set to 3 requests per second and parallel processing is used
func (n *videoGamePricesNotifier) updateNotionWishlistItems(
	ctx context.Context,
	convertedNWishList map[model.SteamAppID]*model.NotionWishlistItem,
	listToUpdate map[model.SteamAppID]*model.SteamStoreVideoGameDetails,
) (map[model.SteamAppID]*model.DiscordContent, error) {
	discordContents := make(map[model.SteamAppID]*model.DiscordContent, 0)
	limiter := rate.NewLimiter(3, 1)
	meg := &multierror.Group{}
	for i, v := range listToUpdate {
		if err := limiter.Wait(ctx); err != nil {
			slog.ErrorContext(ctx, "failed to wait the rate limiter", slog.Any("error", err))
			return nil, err
		}

		meg.Go(func() error {
			// Convert the current price of a video game to uint64
			currentPrice, err := n.convertCurrentPrice(ctx, v.CurrentPrice)
			if err != nil {
				return err
			}

			// Compare the current price of a video game with its lowest price
			var lowestPrice *uint64
			if convertedNWishList[i].Properties.LowestPrice.Number == nil || currentPrice == nil {
				// The current and lowest prices are set to nil if either price is not available
				lowestPrice = nil
			} else if *convertedNWishList[i].Properties.LowestPrice.Number > *currentPrice {
				discordContents[i] = &model.DiscordContent{
					Title:        v.Title,
					CurrentPrice: *currentPrice,
					LowestPrice:  *convertedNWishList[i].Properties.LowestPrice.Number,
				}
				lowestPrice = currentPrice
			} else {
				lowestPrice = convertedNWishList[i].Properties.LowestPrice.Number
			}

			input := &service.UpdateNotionWishlistItemInput{
				WishlistItem: &model.NotionWishlistItem{
					ID: convertedNWishList[i].ID,
					Properties: &model.NotionProperties{
						NotionAppID: &model.NotionAppID{
							Title: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: strconv.Itoa(int(i)),
									},
								},
							},
						},
						NotionName: &model.NotionName{
							RichText: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: v.Title,
									},
								},
							},
						},
						CurrentPrice: &model.NotionPrice{
							Number: currentPrice,
						},
						LowestPrice: &model.NotionPrice{
							Number: lowestPrice,
						},
						NotionReleaseDate: &model.NotionReleaseDate{
							NotionDate: n.convertReleaseDate(ctx, v.ReleaseDate),
						},
					},
				},
			}
			if _, err := n.nWIUpdater.UpdateNotionWishlistItem(ctx, input); err != nil {
				slog.ErrorContext(ctx, "failed to update a wishlist item on the Notion DB", slog.Any("error", err))
				return err
			}

			return nil
		})
	}

	if err := meg.Wait(); err != nil {
		slog.ErrorContext(ctx, "failed to update a wishlist item on the Notion DB", slog.Any("error", err))
		return nil, err
	}

	return discordContents, nil
}

// Convert the current price of a video game to uint64
func (n *videoGamePricesNotifier) convertCurrentPrice(
	ctx context.Context,
	currentPrice *model.SteamCurrentPrice,
) (*uint64, error) {
	if currentPrice == nil {
		return nil, nil
	}

	convertedPrice, err := currentPrice.ConvertPriceFormat(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to convert the current price to uint64", slog.Any("error", err))
		return nil, err
	}

	return convertedPrice, nil
}

// Convert the release date of a video game to "2 Jan, 2006" format time string
// [FYI]
// The release date varies depending on the video game
// e.g. "1 Nov, 2024", "2025", and "To be announced"
func (n *videoGamePricesNotifier) convertReleaseDate(
	ctx context.Context,
	releaseDate *model.SteamReleaseDate,
) *model.NotionDate {
	convertedDate, err := releaseDate.ToTime(ctx)
	if releaseDate.Date == "To be announced" || err != nil {
		slog.WarnContext(ctx, "failed to convert the release date to time.Time", slog.Any("error", err))
		return nil
	}

	return &model.NotionDate{
		Start: convertedDate.Format(time.DateOnly),
	}
}

// Delete a wishlist on the Notion DB
//
// [FYI]
// The rate limiter is set to 3 requests per second and parallel processing is used
func (n *videoGamePricesNotifier) deleteNotionWishlistItems(
	ctx context.Context,
	vGDList map[model.SteamAppID]*model.SteamStoreVideoGameDetails,
	nWishList []*model.NotionWishlistItem,
) error {
	// Convert a Notion wishlist to a map
	convertedNWishList := make(map[model.SteamAppID]*model.NotionWishlistItem, len(nWishList))
	for _, v := range nWishList {
		appID, err := strconv.Atoi(v.Properties.NotionAppID.Title[0].NotionText.NotionContent)
		if err != nil {
			slog.ErrorContext(ctx, "failed to convert the app ID to int", slog.Any("error", err))
			return err
		}
		convertedNWishList[model.SteamAppID(appID)] = v
	}

	// Categorize the Notion wishlist items to delete
	listToDelete := make(map[model.SteamAppID]*model.NotionWishlistItem, len(convertedNWishList))
	for i, v := range convertedNWishList {
		if _, ok := vGDList[i]; !ok {
			listToDelete[i] = v
		}
	}

	limiter := rate.NewLimiter(3, 1)
	meg := &multierror.Group{}
	for _, v := range listToDelete {
		if err := limiter.Wait(ctx); err != nil {
			slog.ErrorContext(ctx, "failed to wait the rate limiter", slog.Any("error", err))
			return err
		}

		meg.Go(func() error {
			input := &service.DeleteNotionWishlistItemInput{
				WishlistItem: v,
			}
			if _, err := n.nWIDeleter.DeleteNotionWishlistItem(ctx, input); err != nil {
				slog.ErrorContext(ctx, "failed to delete a wishlist item on the Notion DB", slog.Any("error", err))
				return err
			}

			return nil
		})
	}

	if err := meg.Wait(); err != nil {
		slog.ErrorContext(ctx, "failed to delete a wishlist item on the Notion DB", slog.Any("error", err))
		return err
	}

	return nil
}

type errorOnDiscordNotifier struct {
	cfg         *config.DiscordConfig
	eODNotifier service.ErrorOnDiscordNotifier
}

var _ usecase.ErrorNotifier = (*errorOnDiscordNotifier)(nil)

// Generate a new errorOnDiscordNotifier
func NewErrorOnDiscordNotifier(
	cfg *config.DiscordConfig,
	eODNotifier service.ErrorOnDiscordNotifier,
) *errorOnDiscordNotifier {
	return &errorOnDiscordNotifier{
		cfg:         cfg,
		eODNotifier: eODNotifier,
	}
}

// Notify an error on Discord
func (n *errorOnDiscordNotifier) NotifyError(
	ctx context.Context,
	input *usecase.NotifyErrorInput,
) (*usecase.NotifyErrorOutput, error) {
	eODInput := &service.NotifyErrorOnDiscordInput{
		GeneratedError: input.GeneratedError,
	}
	if _, err := n.eODNotifier.NotifyErrorOnDiscord(ctx, eODInput); err != nil {
		slog.ErrorContext(ctx, "failed to notify an error on Discord", slog.Any("error", err))
		return nil, err
	}

	return &usecase.NotifyErrorOutput{}, nil
}
