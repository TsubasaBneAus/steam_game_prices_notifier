package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/TsubasaBneAus/steam_game_price_notifier/app/model"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/service"
	"github.com/TsubasaBneAus/steam_game_price_notifier/config"
)

const notionAPIURL string = "https://api.notion.com/v1"

type notionWishlistGetter struct {
	cfg        *config.NotionConfig
	httpClient service.HTTPClient
}

var _ service.NotionWishlistGetter = (*notionWishlistGetter)(nil)

// Generate a new NotionWishlistGetter
func NewNotionWishlistGetter(
	cfg *config.NotionConfig,
	httpClient service.HTTPClient,
) *notionWishlistGetter {
	return &notionWishlistGetter{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

// Get a wishlist from the Notion DB
//
// [FYI]
// Pagination is implemented in the API, and it returns 100 items at a time as a default.
// A request is sent with a start cursor, the API returns 100 items + a next cursor to get the next page.
// If the start cursor is nil, the API returns the first page, so the first request is sent without a start cursor.
// If the next cursor is nil, it means that there are no more items to get.
// ref: https://developers.notion.com/reference/post-database-query
func (g *notionWishlistGetter) GetNotionWishlist(
	ctx context.Context,
	input *service.GetNotionWishlistInput,
) (*service.GetNotionWishlistOutput, error) {
	var startCursor *string
	allWishlistItems := make([]*model.NotionWishlistItem, 0)
	for {
		wishlistItems, err := g.getNotionWishlist(ctx, startCursor)
		if err != nil {
			slog.ErrorContext(ctx, "failed to get a wishlist from the Notion DB", slog.Any("error", err))
			return nil, err
		}

		allWishlistItems = append(allWishlistItems, wishlistItems.Results...)

		if wishlistItems.NextCursor == nil {
			break
		}

		startCursor = wishlistItems.NextCursor
	}

	return &service.GetNotionWishlistOutput{
		WishlistItems: allWishlistItems,
	}, nil
}

// Get a wishlist from the Notion DB with a start cursor
func (g *notionWishlistGetter) getNotionWishlist(
	ctx context.Context,
	startCursor *string,
) (*model.NotionWishlistItems, error) {
	reqURL, err := url.JoinPath(notionAPIURL, "databases", g.cfg.NotionDatabaseID, "query")
	if err != nil {
		slog.ErrorContext(ctx, "failed to build a Notion API URL", slog.Any("error", err))
		return nil, err
	}

	body := &model.NotionDBQuerierBody{
		StartCursor: startCursor,
	}
	reqJSON, err := json.Marshal(body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal a Notion API request body", slog.Any("error", err))
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(reqJSON))
	if err != nil {
		slog.ErrorContext(ctx, "failed to create a Notion API request", slog.Any("error", err))
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.cfg.NotionAPIKey))
	req.Header.Set("Notion-Version", "2022-06-28")

	res, err := g.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send a Notion API request", slog.Any("error", err))
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "unexpected status code in a Notion API response", slog.Any("status_code", res.StatusCode))
		return nil, errUnexpectedStatusCode
	}

	wishlistItems := &model.NotionWishlistItems{}
	if err := json.NewDecoder(res.Body).Decode(wishlistItems); err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal a Notion API response", slog.Any("error", err))
		return nil, err
	}

	return wishlistItems, nil
}

type notionWishlistItemCreator struct {
	cfg        *config.NotionConfig
	httpClient service.HTTPClient
}

var _ service.NotionWishlistItemCreator = (*notionWishlistItemCreator)(nil)

// Generate a new NotionWishlistItemCreator
func NewNotionWishlistItemCreator(
	cfg *config.NotionConfig,
	httpClient service.HTTPClient,
) *notionWishlistItemCreator {
	return &notionWishlistItemCreator{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

// Create a wishlist item in the Notion DB
// ref. https://developers.notion.com/reference/post-page
func (c *notionWishlistItemCreator) CreateNotionWishlistItem(
	ctx context.Context,
	input *service.CreateNotionWishlistItemInput,
) (*service.CreateNotionWishlistItemOutput, error) {
	reqURL, err := url.JoinPath(notionAPIURL, "pages")
	if err != nil {
		slog.ErrorContext(ctx, "failed to build a Notion API URL", slog.Any("error", err))
		return nil, err
	}

	reqJSON, err := json.Marshal(input.WishlistItem)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal a Notion API request body", slog.Any("error", err))
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(reqJSON))
	if err != nil {
		slog.ErrorContext(ctx, "failed to create a Notion API request", slog.Any("error", err))
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.cfg.NotionAPIKey))
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send a Notion API request", slog.Any("error", err))
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.ErrorContext(
			ctx,
			"unexpected status code in a Notion API response",
			slog.Any("status_code", res.StatusCode),
		)
		return nil, errUnexpectedStatusCode
	}

	return &service.CreateNotionWishlistItemOutput{}, nil
}

type notionWishlistItemUpdater struct {
	cfg        *config.NotionConfig
	httpClient service.HTTPClient
}

var _ service.NotionWishlistItemUpdater = (*notionWishlistItemUpdater)(nil)

// Generate a new NotionWishlistItemUpdater
func NewNotionWishlistItemUpdater(
	cfg *config.NotionConfig,
	httpClient service.HTTPClient,
) *notionWishlistItemUpdater {
	return &notionWishlistItemUpdater{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

// Update a wishlist item in Notion DB
// ref. https://developers.notion.com/reference/patch-page
func (u *notionWishlistItemUpdater) UpdateNotionWishlistItem(
	ctx context.Context,
	input *service.UpdateNotionWishlistItemInput,
) (*service.UpdateNotionWishlistItemOutput, error) {
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, reqURL, bytes.NewBuffer(reqJSON))
	if err != nil {
		slog.ErrorContext(ctx, "failed to create a Notion API request", slog.Any("error", err))
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", u.cfg.NotionAPIKey))
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	res, err := u.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send a Notion API request", slog.Any("error", err))
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.ErrorContext(
			ctx,
			"unexpected status code in a Notion API response",
			slog.Any("status_code", res.StatusCode),
		)
		return nil, errUnexpectedStatusCode
	}

	return &service.UpdateNotionWishlistItemOutput{}, nil
}

type notionWishlistItemDeleter struct {
	cfg        *config.NotionConfig
	httpClient service.HTTPClient
}

var _ service.NotionWishlistItemDeleter = (*notionWishlistItemDeleter)(nil)

// Generate a new NotionWishlistItemDeleter
func NewNotionWishlistItemDeleter(
	cfg *config.NotionConfig,
	httpClient service.HTTPClient,
) *notionWishlistItemDeleter {
	return &notionWishlistItemDeleter{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

// Delete a wishlist item from the Notion DB
// ref. https://developers.notion.com/reference/archive-a-page
func (d *notionWishlistItemDeleter) DeleteNotionWishlistItem(
	ctx context.Context,
	input *service.DeleteNotionWishlistItemInput,
) (*service.DeleteNotionWishlistItemOutput, error) {
	reqURL, err := url.JoinPath(notionAPIURL, "pages", string(input.WishlistItem.ID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to build a Notion API URL", slog.Any("error", err))
		return nil, err
	}

	reqBody := bytes.NewBufferString(`{ "in_trash": true }`)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, reqURL, reqBody)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create a Notion API request", slog.Any("error", err))
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.cfg.NotionAPIKey))
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	res, err := d.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send a Notion API request", slog.Any("error", err))
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.ErrorContext(
			ctx,
			"unexpected status code in a Notion API response",
			slog.Any("status_code", res.StatusCode),
		)
		return nil, errUnexpectedStatusCode
	}

	return &service.DeleteNotionWishlistItemOutput{}, nil
}
