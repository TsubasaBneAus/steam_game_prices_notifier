package model

import (
	"context"
	"log/slog"
	"time"
)

// A Notion page ID
type NotionPageID string

// Items of a wishlist in the Notion DB
type NotionWishlistItems struct {
	Results []*NotionWishlistItem `json:"results"`
}

// An item of NotionWishlistItems
type NotionWishlistItem struct {
	ID         NotionPageID      `json:"id"`
	Parent     *NotionParent     `json:"parent,omitempty"`
	Properties *NotionProperties `json:"properties"`
}

// A Notion database ID
type NotionDatabaseID string

// A parent of NotionWishlistItem
type NotionParent struct {
	DatabaseID NotionDatabaseID `json:"database_id"`
}

// Properties of NotionWishlistItem
type NotionProperties struct {
	NotionAppID       *NotionAppID       `json:"App ID,omitempty"`
	NotionName        *NotionName        `json:"Name,omitempty"`
	CurrentPrice      *NotionPrice       `json:"Current Price,omitempty"`
	LowestPrice       *NotionPrice       `json:"Lowest Price,omitempty"`
	NotionReleaseDate *NotionReleaseDate `json:"Release Date,omitempty"`
}

// An app ID of NotionProperties
type NotionAppID struct {
	Title []*NotionContent `json:"title"`
}

// A name of NotionProperties
type NotionName struct {
	RichText []*NotionContent `json:"rich_text"`
}

// A content of NotionName
type NotionContent struct {
	NotionText *NotionText `json:"text"`
}

// A text of NotionContent
type NotionText struct {
	NotionContent string `json:"content"`
}

// A price of NotionProperties
type NotionPrice struct {
	Number *uint64 `json:"number"`
}

// A release date of NotionProperties
type NotionReleaseDate struct {
	NotionDate *NotionDate `json:"date"`
}

// A date of NotionReleaseDate
type NotionDate struct {
	Start string `json:"start"`
}

// Convert date string to time.Time (JST)
func (d *NotionDate) ToTime(ctx context.Context) (*time.Time, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		slog.ErrorContext(ctx, "failed to load location", slog.Any("error", err))
		return nil, err
	}
	parsedTime, err := time.ParseInLocation(time.DateOnly, d.Start, loc)
	if err != nil {
		slog.ErrorContext(ctx, "failed to parse string to time.Time", slog.Any("error", err))
		return nil, err
	}

	return &parsedTime, nil
}
