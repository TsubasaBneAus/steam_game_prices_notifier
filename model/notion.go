package model

import (
	"context"
	"log/slog"
	"time"
)

type PageID string

type WishlistItems struct {
	Results []*WishlistItem `json:"results"`
}

type WishlistItem struct {
	ID         PageID      `json:"id"`
	Properties *Properties `json:"properties"`
}

type Properties struct {
	AppID        *AppID       `json:"App ID,omitempty"`
	Name         *Name        `json:"Name,omitempty"`
	CurrentPrice *Price       `json:"Current Price,omitempty"`
	LowestPrice  *Price       `json:"Lowest Price,omitempty"`
	ReleaseDate  *ReleaseDate `json:"Release Date,omitempty"`
}

type AppID struct {
	Title []*Content `json:"title"`
}

type Name struct {
	RichText []*Content `json:"rich_text"`
}

type Content struct {
	Text *Text `json:"text"`
}

type Text struct {
	Content string `json:"content"`
}

type Price struct {
	Number *uint64 `json:"number"`
}

type ReleaseDate struct {
	Date *Date `json:"date"`
}

type Date struct {
	Start string `json:"start"`
}

// Convert string to time.Time (JST)
func (d *Date) ToTime(ctx context.Context) (*time.Time, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		slog.ErrorContext(ctx, "failed to load location", slog.Any("error", err))
		return nil, err
	}
	parsedTime, err := time.ParseInLocation(time.DateOnly, "2013-11-03", loc)
	if err != nil {
		slog.ErrorContext(ctx, "failed to parse string to time.Time", slog.Any("error", err))
		return nil, err
	}

	return &parsedTime, nil
}
