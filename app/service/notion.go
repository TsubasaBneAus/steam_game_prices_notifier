package service

import (
	"context"

	"github.com/TsubasaBneAus/steam_game_price_notifier/app/model"
)

//go:generate mockgen -source=./notion.go -destination=../external/notion/mock/notion.go -package=mock -typed

type (
	// An input to get a wishlist from the Notion DB
	GetNotionWishlistInput struct{}

	// An output to get a wishlist from the Notion DB
	GetNotionWishlistOutput struct {
		WishlistItems *model.NotionWishlistItems
	}

	// An interface to get a wishlist from the Notion DB
	NotionWishlistGetter interface {
		GetNotionWishlist(
			ctx context.Context,
			input *GetNotionWishlistInput,
		) (*GetNotionWishlistOutput, error)
	}
)

type (
	// An input to create a wishlist item to the Notion DB
	CreateNotionWishlistItemInput struct {
		WishlistItem *model.NotionWishlistItem
	}

	// An output to create a wishlist item to the Notion DB
	CreateNotionWishlistItemOutput struct{}

	// An interface to create a wishlist item to the Notion DB
	NotionWishlistItemCreator interface {
		CreateNotionWishlistItem(
			ctx context.Context,
			input *CreateNotionWishlistItemInput,
		) (*CreateNotionWishlistItemOutput, error)
	}
)

type (
	// An input to update a wishlist item in the Notion DB
	UpdateNotionWishlistItemInput struct {
		WishlistItem *model.NotionWishlistItem
	}

	// An output to update a wishlist item in the Notion DB
	UpdateNotionWishlistItemOutput struct{}

	// An interface to update a wishlist in the Notion DB
	NotionWishlistItemUpdater interface {
		UpdateNotionWishlistItem(
			ctx context.Context,
			input *UpdateNotionWishlistItemInput,
		) (*UpdateNotionWishlistItemOutput, error)
	}
)

type (
	// An input to delete a wishlist item from the Notion DB
	DeleteNotionWishlistItemInput struct {
		WishlistItem *model.NotionWishlistItem
	}

	// An output to delete a wishlist item from the Notion DB
	DeleteNotionWishlistItemOutput struct{}

	// An interface to delete a wishlist item from the Notion DB
	NotionWishlistItemDeleter interface {
		DeleteNotionWishlistItem(
			ctx context.Context,
			input *DeleteNotionWishlistItemInput,
		) (*DeleteNotionWishlistItemOutput, error)
	}
)
