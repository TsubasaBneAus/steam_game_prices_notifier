package interactor

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	discord "github.com/TsubasaBneAus/steam_game_price_notifier/app/external/discord/mock"
	notion "github.com/TsubasaBneAus/steam_game_price_notifier/app/external/notion/mock"
	steam "github.com/TsubasaBneAus/steam_game_price_notifier/app/external/steam/mock"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/model"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/service"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/usecase"
	"github.com/TsubasaBneAus/steam_game_price_notifier/config"
	"github.com/shogo82148/pointer"
	"go.uber.org/mock/gomock"
)

func TestNotifyVideoGamePrices(t *testing.T) {
	t.Parallel()

	// There is only one record in the Notion DB ([1, Title1, 2000, 1500, 2021-01-01])
	// The Steam wishlist has two records ([1, 2])
	// The Steam video game details has two records ([1, Title1, 1000, 1500, 2021-01-01], [2, Title2, nil, nil, To be announced])
	// A new record will be created in the Notion DB ([2, Title2, nil, nil, nil]))
	// The existing record will be updated ([1, Title1, 1000, 1000, 2021-01-01])
	// {1: {Title1, 1000, 1500}} will be notified on Discord
	t.Run("Positive case: Successfully notify video game prices", func(t *testing.T) {
		t.Parallel()

		// Create mocks
		ctrl := gomock.NewController(t)
		sWGetter := steam.NewMockSteamWishlistGetter(ctrl)
		sVGGetter := steam.NewMockSteamVideoGameDetailsGetter(ctrl)
		nWGetter := notion.NewMockNotionWishlistGetter(ctrl)
		nWICreator := notion.NewMockNotionWishlistItemCreator(ctrl)
		nWIUpdater := notion.NewMockNotionWishlistItemUpdater(ctrl)
		vGPODNotifier := discord.NewMockVideoGamePricesOnDiscordNotifier(ctrl)
		{
			input := &service.GetSteamWishlistInput{}
			output := &service.GetSteamWishlistOutput{
				Wishlist: &model.SteamStoreWishlist{
					Response: &model.SteamStoreResponse{
						Items: []*model.SteamStoreItem{
							{
								AppID: 1,
							},
							{
								AppID: 2,
							},
						},
					},
				},
			}
			sWGetter.EXPECT().GetSteamWishlist(gomock.Any(), input).Return(output, nil)
		}
		{
			input1 := &service.GetSteamVideoGameDetailsInput{
				AppID: 1,
			}
			input2 := &service.GetSteamVideoGameDetailsInput{
				AppID: 2,
			}
			output1 := &service.GetSteamVideoGameDetailsOutput{
				VideoGameDetails: &model.SteamStoreVideoGameDetails{
					AppID: 1,
					Title: "Title1",
					CurrentPrice: &model.SteamCurrentPrice{
						Number: json.Number("100000"),
					},
					ReleaseDate: &model.SteamReleaseDate{
						Date: "01 Jan, 2021",
					},
				},
			}
			output2 := &service.GetSteamVideoGameDetailsOutput{
				VideoGameDetails: &model.SteamStoreVideoGameDetails{
					AppID:        2,
					Title:        "Title2",
					CurrentPrice: nil,
					ReleaseDate: &model.SteamReleaseDate{
						Date: "To be announced",
					},
				},
			}
			sVGGetter.EXPECT().GetSteamVideoGameDetails(gomock.Any(), input1).Return(output1, nil)
			sVGGetter.EXPECT().GetSteamVideoGameDetails(gomock.Any(), input2).Return(output2, nil)
		}
		{
			input := &service.GetNotionWishlistInput{}
			output := &service.GetNotionWishlistOutput{
				WishlistItems: &model.NotionWishlistItems{
					Results: []*model.NotionWishlistItem{
						{
							ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
							Parent: &model.NotionParent{
								DatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
							},
							Properties: &model.NotionProperties{
								NotionAppID: &model.NotionAppID{
									Title: []*model.NotionContent{
										{
											NotionText: &model.NotionText{
												NotionContent: "1",
											},
										},
									},
								},
								NotionName: &model.NotionName{
									RichText: []*model.NotionContent{
										{
											NotionText: &model.NotionText{
												NotionContent: "Title1",
											},
										},
									},
								},
								CurrentPrice: &model.NotionPrice{
									Number: pointer.Ptr(uint64(2000)),
								},
								LowestPrice: &model.NotionPrice{
									Number: pointer.Ptr(uint64(1500)),
								},
								NotionReleaseDate: &model.NotionReleaseDate{
									NotionDate: &model.NotionDate{
										Start: "2021-01-01",
									},
								},
							},
						},
					},
				},
			}
			nWGetter.EXPECT().GetNotionWishlist(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.CreateNotionWishlistItemInput{
				WishlistItem: &model.NotionWishlistItem{
					Parent: &model.NotionParent{
						DatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					},
					Properties: &model.NotionProperties{
						NotionAppID: &model.NotionAppID{
							Title: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: "2",
									},
								},
							},
						},
						NotionName: &model.NotionName{
							RichText: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: "Title2",
									},
								},
							},
						},
						CurrentPrice: &model.NotionPrice{
							Number: nil,
						},
						LowestPrice: &model.NotionPrice{
							Number: nil,
						},
						NotionReleaseDate: &model.NotionReleaseDate{
							NotionDate: nil,
						},
					},
				},
			}
			output := &service.CreateNotionWishlistItemOutput{}
			nWICreator.EXPECT().CreateNotionWishlistItem(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.UpdateNotionWishlistItemInput{
				WishlistItem: &model.NotionWishlistItem{
					ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					Properties: &model.NotionProperties{
						NotionAppID: &model.NotionAppID{
							Title: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: "1",
									},
								},
							},
						},
						NotionName: &model.NotionName{
							RichText: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: "Title1",
									},
								},
							},
						},
						CurrentPrice: &model.NotionPrice{
							Number: pointer.Ptr(uint64(1000)),
						},
						LowestPrice: &model.NotionPrice{
							Number: pointer.Ptr(uint64(1000)),
						},
						NotionReleaseDate: &model.NotionReleaseDate{
							NotionDate: &model.NotionDate{
								Start: "2021-01-01",
							},
						},
					},
				},
			}
			output := &service.UpdateNotionWishlistItemOutput{}
			nWIUpdater.EXPECT().UpdateNotionWishlistItem(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.NotifyVideoGamePricesOnDiscordInput{
				DiscordContents: map[model.SteamAppID]*model.DiscordContent{
					1: {
						Title:        "Title1",
						CurrentPrice: 1000,
						LowestPrice:  1500,
					},
				},
			}
			output := &service.NotifyVideoGamePricesOnDiscordOutput{}
			vGPODNotifier.EXPECT().NotifyVideoGamePricesOnDiscord(gomock.Any(), input).Return(output, nil)
		}

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy-notion-api-key",
			NotionDatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		}
		n := NewGamePricesNotifier(cfg, sWGetter, sVGGetter, nWGetter, nWICreator, nWIUpdater, vGPODNotifier)
		input := &usecase.NotifyVideoGamePricesInput{}
		if _, err := n.NotifyVideoGamePrices(ctx, input); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Positive case: The lowest price of a certain record is not filled in the Notion DB", func(t *testing.T) {
		t.Parallel()

		// Create mocks
		ctrl := gomock.NewController(t)
		sWGetter := steam.NewMockSteamWishlistGetter(ctrl)
		sVGGetter := steam.NewMockSteamVideoGameDetailsGetter(ctrl)
		nWGetter := notion.NewMockNotionWishlistGetter(ctrl)
		nWIUpdater := notion.NewMockNotionWishlistItemUpdater(ctrl)
		{
			input := &service.GetSteamWishlistInput{}
			output := &service.GetSteamWishlistOutput{
				Wishlist: &model.SteamStoreWishlist{
					Response: &model.SteamStoreResponse{
						Items: []*model.SteamStoreItem{
							{
								AppID: 1,
							},
						},
					},
				},
			}
			sWGetter.EXPECT().GetSteamWishlist(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.GetSteamVideoGameDetailsInput{
				AppID: 1,
			}
			output := &service.GetSteamVideoGameDetailsOutput{
				VideoGameDetails: &model.SteamStoreVideoGameDetails{
					AppID: 1,
					Title: "Title1",
					CurrentPrice: &model.SteamCurrentPrice{
						Number: json.Number("100000"),
					},
					ReleaseDate: &model.SteamReleaseDate{
						Date: "01 Jan, 2021",
					},
				},
			}
			sVGGetter.EXPECT().GetSteamVideoGameDetails(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.GetNotionWishlistInput{}
			output := &service.GetNotionWishlistOutput{
				WishlistItems: &model.NotionWishlistItems{
					Results: []*model.NotionWishlistItem{
						{
							ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
							Parent: &model.NotionParent{
								DatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
							},
							Properties: &model.NotionProperties{
								NotionAppID: &model.NotionAppID{
									Title: []*model.NotionContent{
										{
											NotionText: &model.NotionText{
												NotionContent: "1",
											},
										},
									},
								},
								NotionName: &model.NotionName{
									RichText: []*model.NotionContent{
										{
											NotionText: &model.NotionText{
												NotionContent: "Title1",
											},
										},
									},
								},
								CurrentPrice: &model.NotionPrice{
									Number: pointer.Ptr(uint64(2000)),
								},
								LowestPrice: &model.NotionPrice{
									Number: nil,
								},
								NotionReleaseDate: &model.NotionReleaseDate{
									NotionDate: &model.NotionDate{
										Start: "2021-01-01",
									},
								},
							},
						},
					},
				},
			}
			nWGetter.EXPECT().GetNotionWishlist(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.UpdateNotionWishlistItemInput{
				WishlistItem: &model.NotionWishlistItem{
					ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					Properties: &model.NotionProperties{
						NotionAppID: &model.NotionAppID{
							Title: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: "1",
									},
								},
							},
						},
						NotionName: &model.NotionName{
							RichText: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: "Title1",
									},
								},
							},
						},
						CurrentPrice: &model.NotionPrice{
							Number: pointer.Ptr(uint64(1000)),
						},
						LowestPrice: &model.NotionPrice{
							Number: nil,
						},
						NotionReleaseDate: &model.NotionReleaseDate{
							NotionDate: &model.NotionDate{
								Start: "2021-01-01",
							},
						},
					},
				},
			}
			output := &service.UpdateNotionWishlistItemOutput{}
			nWIUpdater.EXPECT().UpdateNotionWishlistItem(gomock.Any(), input).Return(output, nil)
		}

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy-notion-api-key",
			NotionDatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		}
		n := NewGamePricesNotifier(cfg, sWGetter, sVGGetter, nWGetter, nil, nWIUpdater, nil)
		input := &usecase.NotifyVideoGamePricesInput{}
		if _, err := n.NotifyVideoGamePrices(ctx, input); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Positive case: The current price of a video game is more expensive than its lowest price", func(t *testing.T) {
		t.Parallel()

		// Create mocks
		ctrl := gomock.NewController(t)
		sWGetter := steam.NewMockSteamWishlistGetter(ctrl)
		sVGGetter := steam.NewMockSteamVideoGameDetailsGetter(ctrl)
		nWGetter := notion.NewMockNotionWishlistGetter(ctrl)
		nWIUpdater := notion.NewMockNotionWishlistItemUpdater(ctrl)
		{
			input := &service.GetSteamWishlistInput{}
			output := &service.GetSteamWishlistOutput{
				Wishlist: &model.SteamStoreWishlist{
					Response: &model.SteamStoreResponse{
						Items: []*model.SteamStoreItem{
							{
								AppID: 1,
							},
						},
					},
				},
			}
			sWGetter.EXPECT().GetSteamWishlist(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.GetSteamVideoGameDetailsInput{
				AppID: 1,
			}
			output := &service.GetSteamVideoGameDetailsOutput{
				VideoGameDetails: &model.SteamStoreVideoGameDetails{
					AppID: 1,
					Title: "Title1",
					CurrentPrice: &model.SteamCurrentPrice{
						Number: json.Number("200000"),
					},
					ReleaseDate: &model.SteamReleaseDate{
						Date: "01 Jan, 2021",
					},
				},
			}
			sVGGetter.EXPECT().GetSteamVideoGameDetails(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.GetNotionWishlistInput{}
			output := &service.GetNotionWishlistOutput{
				WishlistItems: &model.NotionWishlistItems{
					Results: []*model.NotionWishlistItem{
						{
							ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
							Parent: &model.NotionParent{
								DatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
							},
							Properties: &model.NotionProperties{
								NotionAppID: &model.NotionAppID{
									Title: []*model.NotionContent{
										{
											NotionText: &model.NotionText{
												NotionContent: "1",
											},
										},
									},
								},
								NotionName: &model.NotionName{
									RichText: []*model.NotionContent{
										{
											NotionText: &model.NotionText{
												NotionContent: "Title1",
											},
										},
									},
								},
								CurrentPrice: &model.NotionPrice{
									Number: pointer.Ptr(uint64(2000)),
								},
								LowestPrice: &model.NotionPrice{
									Number: pointer.Ptr(uint64(1500)),
								},
								NotionReleaseDate: &model.NotionReleaseDate{
									NotionDate: &model.NotionDate{
										Start: "2021-01-01",
									},
								},
							},
						},
					},
				},
			}
			nWGetter.EXPECT().GetNotionWishlist(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.UpdateNotionWishlistItemInput{
				WishlistItem: &model.NotionWishlistItem{
					ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					Properties: &model.NotionProperties{
						NotionAppID: &model.NotionAppID{
							Title: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: "1",
									},
								},
							},
						},
						NotionName: &model.NotionName{
							RichText: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: "Title1",
									},
								},
							},
						},
						CurrentPrice: &model.NotionPrice{
							Number: pointer.Ptr(uint64(2000)),
						},
						LowestPrice: &model.NotionPrice{
							Number: pointer.Ptr(uint64(1500)),
						},
						NotionReleaseDate: &model.NotionReleaseDate{
							NotionDate: &model.NotionDate{
								Start: "2021-01-01",
							},
						},
					},
				},
			}
			output := &service.UpdateNotionWishlistItemOutput{}
			nWIUpdater.EXPECT().UpdateNotionWishlistItem(gomock.Any(), input).Return(output, nil)
		}

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy-notion-api-key",
			NotionDatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		}
		n := NewGamePricesNotifier(cfg, sWGetter, sVGGetter, nWGetter, nil, nWIUpdater, nil)
		input := &usecase.NotifyVideoGamePricesInput{}
		if _, err := n.NotifyVideoGamePrices(ctx, input); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Negative case: Failed to get a Steam Store wishlist", func(t *testing.T) {
		t.Parallel()

		// Create mocks
		ctrl := gomock.NewController(t)
		sWGetter := steam.NewMockSteamWishlistGetter(ctrl)
		wantErr := errors.New("unexpected error")
		{
			input := &service.GetSteamWishlistInput{}
			sWGetter.EXPECT().GetSteamWishlist(gomock.Any(), input).Return(nil, wantErr)
		}

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy-notion-api-key",
			NotionDatabaseID: "dummy-notion-database-id",
		}
		n := NewGamePricesNotifier(cfg, sWGetter, nil, nil, nil, nil, nil)
		input := &usecase.NotifyVideoGamePricesInput{}
		if _, gotErr := n.NotifyVideoGamePrices(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Failed to get a Steam Store video game details", func(t *testing.T) {
		t.Parallel()

		// Create mocks
		ctrl := gomock.NewController(t)
		sWGetter := steam.NewMockSteamWishlistGetter(ctrl)
		sVGGetter := steam.NewMockSteamVideoGameDetailsGetter(ctrl)
		wantErr := errors.New("unexpected error")
		{
			input := &service.GetSteamWishlistInput{}
			output := &service.GetSteamWishlistOutput{
				Wishlist: &model.SteamStoreWishlist{
					Response: &model.SteamStoreResponse{
						Items: []*model.SteamStoreItem{
							{
								AppID: 1,
							},
						},
					},
				},
			}
			sWGetter.EXPECT().GetSteamWishlist(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.GetSteamVideoGameDetailsInput{
				AppID: 1,
			}
			sVGGetter.EXPECT().GetSteamVideoGameDetails(gomock.Any(), input).Return(nil, wantErr)
		}

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy-notion-api-key",
			NotionDatabaseID: "dummy-notion-database-id",
		}
		n := NewGamePricesNotifier(cfg, sWGetter, sVGGetter, nil, nil, nil, nil)
		input := &usecase.NotifyVideoGamePricesInput{}
		if _, gotErr := n.NotifyVideoGamePrices(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Failed to get a Notion wishlist", func(t *testing.T) {
		t.Parallel()

		// Create mocks
		ctrl := gomock.NewController(t)
		sWGetter := steam.NewMockSteamWishlistGetter(ctrl)
		sVGGetter := steam.NewMockSteamVideoGameDetailsGetter(ctrl)
		nWGetter := notion.NewMockNotionWishlistGetter(ctrl)
		wantErr := errors.New("unexpected error")
		{
			input := &service.GetSteamWishlistInput{}
			output := &service.GetSteamWishlistOutput{
				Wishlist: &model.SteamStoreWishlist{
					Response: &model.SteamStoreResponse{
						Items: []*model.SteamStoreItem{
							{
								AppID: 1,
							},
						},
					},
				},
			}
			sWGetter.EXPECT().GetSteamWishlist(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.GetSteamVideoGameDetailsInput{
				AppID: 1,
			}
			output := &service.GetSteamVideoGameDetailsOutput{
				VideoGameDetails: &model.SteamStoreVideoGameDetails{
					AppID: 1,
					Title: "Title1",
					CurrentPrice: &model.SteamCurrentPrice{
						Number: json.Number("100000"),
					},
					ReleaseDate: &model.SteamReleaseDate{
						Date: "01 Jan, 2021",
					},
				},
			}
			sVGGetter.EXPECT().GetSteamVideoGameDetails(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.GetNotionWishlistInput{}
			nWGetter.EXPECT().GetNotionWishlist(gomock.Any(), input).Return(nil, wantErr)
		}

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy-notion-api-key",
			NotionDatabaseID: "dummy-notion-database-id",
		}
		n := NewGamePricesNotifier(cfg, sWGetter, sVGGetter, nWGetter, nil, nil, nil)
		input := &usecase.NotifyVideoGamePricesInput{}
		if _, gotErr := n.NotifyVideoGamePrices(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Failed to create a Notion wishlist item", func(t *testing.T) {
		t.Parallel()

		// Create mocks
		ctrl := gomock.NewController(t)
		sWGetter := steam.NewMockSteamWishlistGetter(ctrl)
		sVGGetter := steam.NewMockSteamVideoGameDetailsGetter(ctrl)
		nWGetter := notion.NewMockNotionWishlistGetter(ctrl)
		nWICreator := notion.NewMockNotionWishlistItemCreator(ctrl)
		wantErr := errors.New("unexpected error")
		{
			input := &service.GetSteamWishlistInput{}
			output := &service.GetSteamWishlistOutput{
				Wishlist: &model.SteamStoreWishlist{
					Response: &model.SteamStoreResponse{
						Items: []*model.SteamStoreItem{
							{
								AppID: 1,
							},
						},
					},
				},
			}
			sWGetter.EXPECT().GetSteamWishlist(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.GetSteamVideoGameDetailsInput{
				AppID: 1,
			}
			output := &service.GetSteamVideoGameDetailsOutput{
				VideoGameDetails: &model.SteamStoreVideoGameDetails{
					AppID: 1,
					Title: "Title1",
					CurrentPrice: &model.SteamCurrentPrice{
						Number: json.Number("100000"),
					},
					ReleaseDate: &model.SteamReleaseDate{
						Date: "01 Jan, 2021",
					},
				},
			}
			sVGGetter.EXPECT().GetSteamVideoGameDetails(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.GetNotionWishlistInput{}
			output := &service.GetNotionWishlistOutput{
				WishlistItems: &model.NotionWishlistItems{
					Results: []*model.NotionWishlistItem{},
				},
			}
			nWGetter.EXPECT().GetNotionWishlist(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.CreateNotionWishlistItemInput{
				WishlistItem: &model.NotionWishlistItem{
					Parent: &model.NotionParent{
						DatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					},
					Properties: &model.NotionProperties{
						NotionAppID: &model.NotionAppID{
							Title: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: "1",
									},
								},
							},
						},
						NotionName: &model.NotionName{
							RichText: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: "Title1",
									},
								},
							},
						},
						CurrentPrice: &model.NotionPrice{
							Number: pointer.Ptr(uint64(1000)),
						},
						LowestPrice: &model.NotionPrice{
							Number: nil,
						},
						NotionReleaseDate: &model.NotionReleaseDate{
							NotionDate: &model.NotionDate{
								Start: "2021-01-01",
							},
						},
					},
				},
			}
			nWICreator.EXPECT().CreateNotionWishlistItem(gomock.Any(), input).Return(nil, wantErr)
		}

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy-notion-api-key",
			NotionDatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		}
		n := NewGamePricesNotifier(cfg, sWGetter, sVGGetter, nWGetter, nWICreator, nil, nil)
		input := &usecase.NotifyVideoGamePricesInput{}
		if _, gotErr := n.NotifyVideoGamePrices(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Failed to update a Notion wishlist item", func(t *testing.T) {
		t.Parallel()

		// Create mocks
		ctrl := gomock.NewController(t)
		sWGetter := steam.NewMockSteamWishlistGetter(ctrl)
		sVGGetter := steam.NewMockSteamVideoGameDetailsGetter(ctrl)
		nWGetter := notion.NewMockNotionWishlistGetter(ctrl)
		nWIUpdater := notion.NewMockNotionWishlistItemUpdater(ctrl)
		wantErr := errors.New("unexpected error")
		{
			input := &service.GetSteamWishlistInput{}
			output := &service.GetSteamWishlistOutput{
				Wishlist: &model.SteamStoreWishlist{
					Response: &model.SteamStoreResponse{
						Items: []*model.SteamStoreItem{
							{
								AppID: 1,
							},
						},
					},
				},
			}
			sWGetter.EXPECT().GetSteamWishlist(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.GetSteamVideoGameDetailsInput{
				AppID: 1,
			}
			output := &service.GetSteamVideoGameDetailsOutput{
				VideoGameDetails: &model.SteamStoreVideoGameDetails{
					AppID: 1,
					Title: "Title1",
					CurrentPrice: &model.SteamCurrentPrice{
						Number: json.Number("100000"),
					},
					ReleaseDate: &model.SteamReleaseDate{
						Date: "01 Jan, 2021",
					},
				},
			}
			sVGGetter.EXPECT().GetSteamVideoGameDetails(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.GetNotionWishlistInput{}
			output := &service.GetNotionWishlistOutput{
				WishlistItems: &model.NotionWishlistItems{
					Results: []*model.NotionWishlistItem{
						{
							ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
							Parent: &model.NotionParent{
								DatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
							},
							Properties: &model.NotionProperties{
								NotionAppID: &model.NotionAppID{
									Title: []*model.NotionContent{
										{
											NotionText: &model.NotionText{
												NotionContent: "1",
											},
										},
									},
								},
								NotionName: &model.NotionName{
									RichText: []*model.NotionContent{
										{
											NotionText: &model.NotionText{
												NotionContent: "Title1",
											},
										},
									},
								},
								CurrentPrice: &model.NotionPrice{
									Number: pointer.Ptr(uint64(2000)),
								},
								LowestPrice: &model.NotionPrice{
									Number: pointer.Ptr(uint64(1500)),
								},
								NotionReleaseDate: &model.NotionReleaseDate{
									NotionDate: &model.NotionDate{
										Start: "2021-01-01",
									},
								},
							},
						},
					},
				},
			}
			nWGetter.EXPECT().GetNotionWishlist(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.UpdateNotionWishlistItemInput{
				WishlistItem: &model.NotionWishlistItem{
					ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					Properties: &model.NotionProperties{
						NotionAppID: &model.NotionAppID{
							Title: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: "1",
									},
								},
							},
						},
						NotionName: &model.NotionName{
							RichText: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: "Title1",
									},
								},
							},
						},
						CurrentPrice: &model.NotionPrice{
							Number: pointer.Ptr(uint64(1000)),
						},
						LowestPrice: &model.NotionPrice{
							Number: pointer.Ptr(uint64(1000)),
						},
						NotionReleaseDate: &model.NotionReleaseDate{
							NotionDate: &model.NotionDate{
								Start: "2021-01-01",
							},
						},
					},
				},
			}
			nWIUpdater.EXPECT().UpdateNotionWishlistItem(gomock.Any(), input).Return(nil, wantErr)
		}

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy-notion-api-key",
			NotionDatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		}
		n := NewGamePricesNotifier(cfg, sWGetter, sVGGetter, nWGetter, nil, nWIUpdater, nil)
		input := &usecase.NotifyVideoGamePricesInput{}
		if _, gotErr := n.NotifyVideoGamePrices(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Failed to notify video game prices on Discord", func(t *testing.T) {
		t.Parallel()

		// Create mocks
		ctrl := gomock.NewController(t)
		sWGetter := steam.NewMockSteamWishlistGetter(ctrl)
		sVGGetter := steam.NewMockSteamVideoGameDetailsGetter(ctrl)
		nWGetter := notion.NewMockNotionWishlistGetter(ctrl)
		nWIUpdater := notion.NewMockNotionWishlistItemUpdater(ctrl)
		vGPODNotifier := discord.NewMockVideoGamePricesOnDiscordNotifier(ctrl)
		wantErr := errors.New("unexpected error")
		{
			input := &service.GetSteamWishlistInput{}
			output := &service.GetSteamWishlistOutput{
				Wishlist: &model.SteamStoreWishlist{
					Response: &model.SteamStoreResponse{
						Items: []*model.SteamStoreItem{
							{
								AppID: 1,
							},
						},
					},
				},
			}
			sWGetter.EXPECT().GetSteamWishlist(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.GetSteamVideoGameDetailsInput{
				AppID: 1,
			}
			output := &service.GetSteamVideoGameDetailsOutput{
				VideoGameDetails: &model.SteamStoreVideoGameDetails{
					AppID: 1,
					Title: "Title1",
					CurrentPrice: &model.SteamCurrentPrice{
						Number: json.Number("100000"),
					},
					ReleaseDate: &model.SteamReleaseDate{
						Date: "01 Jan, 2021",
					},
				},
			}
			sVGGetter.EXPECT().GetSteamVideoGameDetails(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.GetNotionWishlistInput{}
			output := &service.GetNotionWishlistOutput{
				WishlistItems: &model.NotionWishlistItems{
					Results: []*model.NotionWishlistItem{
						{
							ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
							Parent: &model.NotionParent{
								DatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
							},
							Properties: &model.NotionProperties{
								NotionAppID: &model.NotionAppID{
									Title: []*model.NotionContent{
										{
											NotionText: &model.NotionText{
												NotionContent: "1",
											},
										},
									},
								},
								NotionName: &model.NotionName{
									RichText: []*model.NotionContent{
										{
											NotionText: &model.NotionText{
												NotionContent: "Title1",
											},
										},
									},
								},
								CurrentPrice: &model.NotionPrice{
									Number: pointer.Ptr(uint64(2000)),
								},
								LowestPrice: &model.NotionPrice{
									Number: pointer.Ptr(uint64(1500)),
								},
								NotionReleaseDate: &model.NotionReleaseDate{
									NotionDate: &model.NotionDate{
										Start: "2021-01-01",
									},
								},
							},
						},
					},
				},
			}
			nWGetter.EXPECT().GetNotionWishlist(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.UpdateNotionWishlistItemInput{
				WishlistItem: &model.NotionWishlistItem{
					ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					Properties: &model.NotionProperties{
						NotionAppID: &model.NotionAppID{
							Title: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: "1",
									},
								},
							},
						},
						NotionName: &model.NotionName{
							RichText: []*model.NotionContent{
								{
									NotionText: &model.NotionText{
										NotionContent: "Title1",
									},
								},
							},
						},
						CurrentPrice: &model.NotionPrice{
							Number: pointer.Ptr(uint64(1000)),
						},
						LowestPrice: &model.NotionPrice{
							Number: pointer.Ptr(uint64(1000)),
						},
						NotionReleaseDate: &model.NotionReleaseDate{
							NotionDate: &model.NotionDate{
								Start: "2021-01-01",
							},
						},
					},
				},
			}
			output := &service.UpdateNotionWishlistItemOutput{}
			nWIUpdater.EXPECT().UpdateNotionWishlistItem(gomock.Any(), input).Return(output, nil)
		}
		{
			input := &service.NotifyVideoGamePricesOnDiscordInput{
				DiscordContents: map[model.SteamAppID]*model.DiscordContent{
					1: {
						Title:        "Title1",
						CurrentPrice: 1000,
						LowestPrice:  1500,
					},
				},
			}
			vGPODNotifier.EXPECT().NotifyVideoGamePricesOnDiscord(gomock.Any(), input).Return(nil, wantErr)
		}

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy-notion-api-key",
			NotionDatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		}
		n := NewGamePricesNotifier(cfg, sWGetter, sVGGetter, nWGetter, nil, nWIUpdater, vGPODNotifier)
		input := &usecase.NotifyVideoGamePricesInput{}
		if _, gotErr := n.NotifyVideoGamePrices(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})
}
