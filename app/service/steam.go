package service

import (
	"context"

	"github.com/TsubasaBneAus/steam_game_price_notifier/app/model"
)

//go:generate mockgen -source=./steam.go -destination=../external/steam/mock/steam.go -package=mock -typed

type (
	// An input to get a wishlist from the Steam Store
	GetSteamWishlistInput struct{}

	// An output to get a wishlist from the Steam Store
	GetSteamWishlistOutput struct {
		Wishlist *model.SteamStoreWishlist
	}

	// An interface to get a wishlist from the Steam Store
	SteamWishlistGetter interface {
		GetSteamWishlist(
			ctx context.Context,
			input *GetSteamWishlistInput,
		) (*GetSteamWishlistOutput, error)
	}
)

type (
	// An input to get video game details from the Steam Store
	GetSteamVideoGameDetailsInput struct {
		AppID model.SteamAppID
	}

	// An output to get video game details from the Steam Store
	GetSteamVideoGameDetailsOutput struct {
		VideoGameDetails *model.SteamStoreVideoGameDetails
	}

	// An interface to get video game details from the Steam Store
	SteamVideoGameDetailsGetter interface {
		GetSteamVideoGameDetails(
			ctx context.Context,
			input *GetSteamVideoGameDetailsInput,
		) (*GetSteamVideoGameDetailsOutput, error)
	}
)
