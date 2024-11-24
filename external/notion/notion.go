package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/TsubasaBneAus/steam_game_price_notifier/config"
	"github.com/TsubasaBneAus/steam_game_price_notifier/external/httpClient"
	"github.com/TsubasaBneAus/steam_game_price_notifier/model"
)

const notionAPIURL string = "https://api.notion.com/v1"

type (
	GetWishlistInput struct{}

	GetWishlistOutput struct {
		WishlistItems *model.WishlistItems
	}

	WishlistGetter interface {
		GetWishlist(ctx context.Context, input *GetWishlistInput) (*GetWishlistOutput, error)
	}
)

type wishlistGetter struct {
	cfg        *config.Envs
	httpClient httpClient.HttpClient
}

var _ WishlistGetter = (*wishlistGetter)(nil)

func NewWishlistGetter(cfg *config.Envs, httpClient httpClient.HttpClient) WishlistGetter {
	return &wishlistGetter{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

// Get Steam Wishlist from Notion database
func (wg *wishlistGetter) GetWishlist(
	ctx context.Context,
	input *GetWishlistInput,
) (*GetWishlistOutput, error) {
	reqURL, err := url.JoinPath(notionAPIURL, "databases", wg.cfg.NotionDatabaseID, "query")
	if err != nil {
		slog.ErrorContext(ctx, "failed to build a Notion API URL", slog.Any("error", err))
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create a Notion API request", slog.Any("error", err))
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", wg.cfg.NotionApiKey))
	req.Header.Set("Notion-Version", "2022-06-28")

	res, err := wg.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send a Notion API request", slog.Any("error", err))
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "unexpected status code in the Notion API response", slog.Any("status_code", res.StatusCode))
		return nil, errUnexpectedStatusCode
	}

	buffer := bytes.Buffer{}
	if _, err := io.Copy(&buffer, res.Body); err != nil {
		slog.ErrorContext(ctx, "failed to read Notion API response", slog.Any("error", err))
		return nil, err
	}

	wishlistItems := &model.WishlistItems{}
	if err := json.Unmarshal(buffer.Bytes(), wishlistItems); err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal a Notion API response", slog.Any("error", err))
		return nil, err
	}

	return &GetWishlistOutput{
		WishlistItems: wishlistItems,
	}, nil
}

type (
	UpdateWishlistInput struct {
		WishlistItem *model.WishlistItem
	}

	UpdateWishlistOutput struct{}

	WishlistUpdater interface {
		UpdateWishlist(ctx context.Context, input *UpdateWishlistInput) (*UpdateWishlistOutput, error)
	}
)

type wishlistUpdater struct {
	cfg        *config.Envs
	httpClient httpClient.HttpClient
}

var _ WishlistUpdater = (*wishlistUpdater)(nil)

func NewWishlistUpdater(cfg *config.Envs, httpClient httpClient.HttpClient) WishlistUpdater {
	return &wishlistUpdater{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

// Update Steam Wishlist in Notion database
func (wu *wishlistUpdater) UpdateWishlist(
	ctx context.Context,
	input *UpdateWishlistInput,
) (*UpdateWishlistOutput, error) {
	reqURL, err := url.JoinPath(notionAPIURL, "pages", string(input.WishlistItem.ID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to build a Notion API URL", slog.Any("error", err))
		return nil, err
	}

	reqJSON, err := json.Marshal(input.WishlistItem)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal a Notion API request body", slog.Any("error", err))
		return nil, err
	}

	fmt.Println(string(reqJSON))

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, reqURL, bytes.NewBuffer(reqJSON))
	if err != nil {
		slog.ErrorContext(ctx, "failed to create a Notion API request", slog.Any("error", err))
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", wu.cfg.NotionApiKey))
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	res, err := wu.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send a Notion API request", slog.Any("error", err))
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "unexpected status code in the Notion API response", slog.Any("status_code", res.StatusCode))
		return nil, errUnexpectedStatusCode
	}

	return &UpdateWishlistOutput{}, nil
}
