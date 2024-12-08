package model

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"
)

// A wishlist on Steam
type SteamStoreWishlist struct {
	Response *SteamStoreResponse `json:"response"`
}

// A response of SteamWishlist
type SteamStoreResponse struct {
	Items []*SteamStoreItem `json:"items"`
}

// An app ID of Steam Store
type SteamAppID uint64

// An Item of SteamResponse
type SteamStoreItem struct {
	AppID uint64 `json:"appid"`
}

// A video game details on Steam
type SteamStoreVideoGameDetails struct {
	AppID        SteamAppID
	Title        string
	CurrentPrice *SteamCurrentPrice
	ReleaseDate  *SteamReleaseDate
}

// A current price of SteamStoreVideoGameDetails
type SteamCurrentPrice struct {
	Number json.Number
}

// Convert the current price format
//
// [FYI]
// Retrieved price contains decimal places
// e.g. 100000 -> 1000 (JPY)
func (p *SteamCurrentPrice) ConvertPriceFormat(ctx context.Context) (*uint64, error) {
	currentPrice, err := p.Number.Int64()
	if err != nil {
		slog.ErrorContext(ctx, "failed to convert the current price to int64", slog.Any("error", err))
		return nil, err
	}

	// Remove the last two digits
	convertedPrice := uint64(currentPrice) / 100

	return &convertedPrice, nil
}

// A release date of SteamStoreVideoGameDetails
type SteamReleaseDate struct {
	Date string
}

// Convert date string into time.Time (JST)
func (d *SteamReleaseDate) ToTime(ctx context.Context) (*time.Time, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		slog.ErrorContext(ctx, "failed to load location", slog.Any("error", err))
		return nil, err
	}
	parsedTime, err := time.ParseInLocation("2 Jan, 2006", d.Date, loc)
	if err != nil {
		slog.WarnContext(ctx, "failed to parse string to time.Time", slog.Any("error", err))
		return nil, err
	}

	return &parsedTime, nil
}
