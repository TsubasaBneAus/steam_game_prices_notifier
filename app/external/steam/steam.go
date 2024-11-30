package steam

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/TsubasaBneAus/steam_game_price_notifier/app/model"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/service"
	"github.com/TsubasaBneAus/steam_game_price_notifier/config"
)

const (
	steamStoreWishlistURL         string = "https://api.steampowered.com/IWishlistService/GetWishlist/v1/"
	steamStoreVideoGameDetailsURL string = "https://store.steampowered.com/api/appdetails/"
)

type steamWishlistGetter struct {
	cfg        *config.SteamConfig
	httpClient service.HTTPClient
}

var _ service.SteamWishlistGetter = (*steamWishlistGetter)(nil)

// Generate a new SteamWishlistGetter
func NewSteamWishlistGetter(
	cfg *config.SteamConfig,
	httpClient service.HTTPClient,
) *steamWishlistGetter {
	return &steamWishlistGetter{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

// Get a wishlist from the Steam Store
func (wg *steamWishlistGetter) GetSteamWishlist(
	ctx context.Context,
	input *service.GetSteamWishlistInput,
) (*service.GetSteamWishlistOutput, error) {
	reqURL, err := url.Parse(steamStoreWishlistURL)
	if err != nil {
		slog.ErrorContext(ctx, "failed to build a Steam Store wishlist URL", slog.Any("error", err))
		return nil, err
	}

	q := reqURL.Query()
	q.Set("steamid", wg.cfg.SteamUserID)
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create a Steam Store wishlist request", slog.Any("error", err))
		return nil, err
	}

	res, err := wg.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send a Steam Store wishlist request", slog.Any("error", err))
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.ErrorContext(
			ctx,
			"unexpected status code in the Steam Store wishlist response",
			slog.Any("status_code", res.StatusCode),
		)
		return nil, errUnexpectedStatusCode
	}

	buffer := bytes.Buffer{}
	if _, err := io.Copy(&buffer, res.Body); err != nil {
		slog.ErrorContext(ctx, "failed to read a Steam Store wishlist response", slog.Any("error", err))
		return nil, err
	}

	wishlist := &model.SteamStoreWishlist{}
	if err := json.Unmarshal(buffer.Bytes(), wishlist); err != nil {
		slog.ErrorContext(
			ctx,
			"failed to unmarshal a Steam Store wishlist response",
			slog.Any("error", err),
		)
		return nil, err
	}

	return &service.GetSteamWishlistOutput{
		Wishlist: wishlist,
	}, nil
}

type steamVideoGameDetailsGetter struct {
	cfg        *config.SteamConfig
	httpClient service.HTTPClient
}

var _ service.SteamVideoGameDetailsGetter = (*steamVideoGameDetailsGetter)(nil)

// Generate a new SteamVideoGameDetailsGetter
func NewSteamVideoGameDetailsGetter(
	cfg *config.SteamConfig,
	httpClient service.HTTPClient,
) *steamVideoGameDetailsGetter {
	return &steamVideoGameDetailsGetter{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

// Get video game details from the Steam Store
func (vg *steamVideoGameDetailsGetter) GetSteamVideoGameDetails(
	ctx context.Context,
	input *service.GetSteamVideoGameDetailsInput,
) (*service.GetSteamVideoGameDetailsOutput, error) {
	reqURL, err := url.Parse(steamStoreVideoGameDetailsURL)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed to build a Steam Store video game details URL",
			slog.Any("error", err),
		)
		return nil, err
	}

	q := reqURL.Query()
	q.Set("cc", "jp")
	q.Set("appids", strconv.FormatUint(uint64(input.AppID), 10))
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed to create a Steam Store video game details request",
			slog.Any("error", err),
		)
		return nil, err
	}

	res, err := vg.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx,
			"failed to send a Steam Store video game details request",
			slog.Any("error", err),
		)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.ErrorContext(
			ctx,
			"unexpected status code in the Steam Store video game details response",
			slog.Any("status_code", res.StatusCode),
		)
		return nil, errUnexpectedStatusCode
	}

	buffer := bytes.Buffer{}
	if _, err := io.Copy(&buffer, res.Body); err != nil {
		slog.ErrorContext(
			ctx,
			"failed to read a Steam Store video game details response",
			slog.Any("error", err),
		)
		return nil, err
	}

	// Unmarshal the response
	//
	// [FYI]
	// Map is used because the response does not have fixed name keys
	videoGameDetails := make(map[string]any, 1)
	decoder := json.NewDecoder(&buffer)
	decoder.UseNumber()
	if err := decoder.Decode(&videoGameDetails); err != nil {
		slog.ErrorContext(
			ctx,
			"failed to unmarshal a Steam Store video game details response",
			slog.Any("error", err),
		)
		return nil, err
	}

	// Extract the data from the response
	appID := videoGameDetails[strconv.FormatUint(uint64(input.AppID), 10)].(map[string]any)
	data := appID["data"].(map[string]any)
	priceOverview := data["price_overview"].(map[string]interface{})
	releaseDate := data["release_date"].(map[string]any)

	return &service.GetSteamVideoGameDetailsOutput{
		VideoGameDetails: &model.SteamStoreVideoGameDetails{
			AppID: input.AppID,
			Title: data["name"].(string),
			CurrentPrice: &model.SteamCurrentPrice{
				Number: priceOverview["final"].(json.Number),
			},
			ReleaseDate: &model.SteamReleaseDate{
				Date: releaseDate["date"].(string),
			},
		},
	}, nil
}
