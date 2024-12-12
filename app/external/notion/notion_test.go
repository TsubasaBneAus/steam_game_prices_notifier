package notion

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	httpclient "github.com/TsubasaBneAus/steam_game_price_notifier/app/external/httpclient/mock"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/model"
	"github.com/TsubasaBneAus/steam_game_price_notifier/app/service"
	"github.com/TsubasaBneAus/steam_game_price_notifier/config"
	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/pointer"
	"go.uber.org/mock/gomock"
)

func TestGetNotionWishlist(t *testing.T) {
	t.Parallel()

	t.Run("Positive case: Successfully get a wishlist from the Notion DB", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(_ *http.Request) (*http.Response, error) {
				jsonFile, err := os.Open("./testdata/wishlist.json")
				if err != nil {
					t.Fatalf("failed to open wishlist.json: %v", err)
				}
				defer jsonFile.Close()

				buffer := bytes.Buffer{}
				if _, err := io.Copy(&buffer, jsonFile); err != nil {
					t.Fatalf("failed to read wishlist.json: %v", err)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(buffer.Bytes())),
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy_notion_api_key",
			NotionDatabaseID: "dummy_notion_database_id",
		}
		wg := NewNotionWishlistGetter(cfg, m)
		got, err := wg.GetNotionWishlist(ctx, &service.GetNotionWishlistInput{})
		if err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
		want := &service.GetNotionWishlistOutput{
			WishlistItems: &model.NotionWishlistItems{
				Results: []*model.NotionWishlistItem{
					{
						ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
						Parent: &model.NotionParent{
							DatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
						},
						Properties: &model.NotionProperties{
							NotionAppID: &model.NotionAppID{
								Title: []*model.NotionContent{},
							},
							NotionName: &model.NotionName{
								RichText: []*model.NotionContent{},
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
					{
						ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
						Parent: &model.NotionParent{
							DatabaseID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
						},
						Properties: &model.NotionProperties{
							NotionAppID: &model.NotionAppID{
								Title: []*model.NotionContent{
									{
										NotionText: &model.NotionText{
											NotionContent: "2701660",
										},
									},
								},
							},
							NotionName: &model.NotionName{
								RichText: []*model.NotionContent{
									{
										NotionText: &model.NotionText{
											NotionContent: "ドラゴンクエストIII　そして伝説へ…",
										},
									},
								},
							},
							CurrentPrice: &model.NotionPrice{
								Number: pointer.Ptr(uint64(7678)),
							},
							LowestPrice: &model.NotionPrice{
								Number: pointer.Ptr(uint64(7678)),
							},
							NotionReleaseDate: &model.NotionReleaseDate{
								NotionDate: &model.NotionDate{
									Start: "2024-11-15",
								},
							},
						},
					},
				},
			},
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("got(-) want(+)\n%s", diff)
		}
	})

	t.Run("Positive case: Successfully Get an empty wishlist from the Notion DB", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(_ *http.Request) (*http.Response, error) {
				jsonFile, err := os.Open("./testdata/empty_wishlist.json")
				if err != nil {
					t.Fatalf("failed to open empty_wishlist.json: %v", err)
				}
				defer jsonFile.Close()

				buffer := bytes.Buffer{}
				if _, err := io.Copy(&buffer, jsonFile); err != nil {
					t.Fatalf("failed to read empty_wishlist.json: %v", err)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(buffer.Bytes())),
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy_notion_api_key",
			NotionDatabaseID: "dummy_notion_database_id",
		}
		wg := NewNotionWishlistGetter(cfg, m)
		got, err := wg.GetNotionWishlist(ctx, &service.GetNotionWishlistInput{})
		if err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
		want := &service.GetNotionWishlistOutput{
			WishlistItems: &model.NotionWishlistItems{
				Results: []*model.NotionWishlistItem{},
			},
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("got(-) want(+)\n%s", diff)
		}
	})

	t.Run("Negative case: Fail to send a request", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		wantErr := errors.New("unexpected error")
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(nil, wantErr)

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy_notion_api_key",
			NotionDatabaseID: "dummy_notion_database_id",
		}
		wg := NewNotionWishlistGetter(cfg, m)
		if _, gotErr := wg.GetNotionWishlist(ctx, &service.GetNotionWishlistInput{}); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Get a status code except 200", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(
				&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       http.NoBody,
				}, nil,
			)

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy_notion_api_key",
			NotionDatabaseID: "dummy_notion_database_id",
		}
		wg := NewNotionWishlistGetter(cfg, m)
		wantErr := errUnexpectedStatusCode
		if _, gotErr := wg.GetNotionWishlist(ctx, &service.GetNotionWishlistInput{}); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Fail to unmarshal a response", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(
				&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("invalid JSON"))),
				}, nil,
			)

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy_notion_api_key",
			NotionDatabaseID: "dummy_notion_database_id",
		}
		wg := NewNotionWishlistGetter(cfg, m)
		if _, err := wg.GetNotionWishlist(ctx, &service.GetNotionWishlistInput{}); err == nil {
			t.Errorf("\ngot: %v\nwant: an error generated by the library", nil)
		}
	})
}

func TestCreateNotionWishlistItem(t *testing.T) {
	t.Parallel()

	t.Run("Positive case: Successfully create a wishlist item in the Notion DB", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(_ *http.Request) (*http.Response, error) {
				jsonFile, err := os.Open("./testdata/created_wishlist_item.json")
				if err != nil {
					t.Fatalf("failed to open created_wishlist_item.json: %v", err)
				}
				defer jsonFile.Close()

				buffer := bytes.Buffer{}
				if _, err := io.Copy(&buffer, jsonFile); err != nil {
					t.Fatalf("failed to read created_wishlist_item.json: %v", err)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(buffer.Bytes())),
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy_notion_api_key",
			NotionDatabaseID: "dummy_notion_database_id",
		}
		wg := NewNotionWishlistItemCreator(cfg, m)
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
									NotionContent: "2701660",
								},
							},
						},
					},
					NotionName: &model.NotionName{
						RichText: []*model.NotionContent{
							{
								NotionText: &model.NotionText{
									NotionContent: "ドラゴンクエストIII　そして伝説へ…",
								},
							},
						},
					},
					CurrentPrice: &model.NotionPrice{
						Number: pointer.Ptr(uint64(7678)),
					},
					LowestPrice: &model.NotionPrice{
						Number: pointer.Ptr(uint64(7678)),
					},
					NotionReleaseDate: &model.NotionReleaseDate{
						NotionDate: &model.NotionDate{
							Start: "2024-11-15",
						},
					},
				},
			},
		}
		if _, err := wg.CreateNotionWishlistItem(ctx, input); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Positive case: Successfully create a wishlist item in the Notion DB with an empty input", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(_ *http.Request) (*http.Response, error) {
				jsonFile, err := os.Open("./testdata/created_empty_wishlist_item.json")
				if err != nil {
					t.Fatalf("failed to open created_empty_wishlist_item.json: %v", err)
				}
				defer jsonFile.Close()

				buffer := bytes.Buffer{}
				if _, err := io.Copy(&buffer, jsonFile); err != nil {
					t.Fatalf("failed to read created_empty_wishlist_item.json: %v", err)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(buffer.Bytes())),
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy_notion_api_key",
			NotionDatabaseID: "dummy_notion_database_id",
		}
		wg := NewNotionWishlistItemCreator(cfg, m)
		input := &service.CreateNotionWishlistItemInput{
			WishlistItem: &model.NotionWishlistItem{
				Parent: &model.NotionParent{
					DatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
				},
				Properties: &model.NotionProperties{
					NotionAppID: &model.NotionAppID{
						Title: nil,
					},
					NotionName: &model.NotionName{
						RichText: []*model.NotionContent{},
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
		if _, err := wg.CreateNotionWishlistItem(ctx, input); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Negative case: Fail to send a request", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		wantErr := errors.New("unexpected error")
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(nil, wantErr)

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy_notion_api_key",
			NotionDatabaseID: "dummy_notion_database_id",
		}
		wg := NewNotionWishlistItemCreator(cfg, m)
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
									NotionContent: "2701660",
								},
							},
						},
					},
					NotionName: &model.NotionName{
						RichText: []*model.NotionContent{
							{
								NotionText: &model.NotionText{
									NotionContent: "ドラゴンクエストIII　そして伝説へ…",
								},
							},
						},
					},
					CurrentPrice: &model.NotionPrice{
						Number: pointer.Ptr(uint64(7678)),
					},
					LowestPrice: &model.NotionPrice{
						Number: pointer.Ptr(uint64(7678)),
					},
					NotionReleaseDate: &model.NotionReleaseDate{
						NotionDate: &model.NotionDate{
							Start: "2024-11-15",
						},
					},
				},
			},
		}
		if _, gotErr := wg.CreateNotionWishlistItem(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Get a status code except 200", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       http.NoBody,
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy_notion_api_key",
			NotionDatabaseID: "dummy_notion_database_id",
		}
		wg := NewNotionWishlistItemCreator(cfg, m)
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
									NotionContent: "2701660",
								},
							},
						},
					},
					NotionName: &model.NotionName{
						RichText: []*model.NotionContent{
							{
								NotionText: &model.NotionText{
									NotionContent: "ドラゴンクエストIII　そして伝説へ…",
								},
							},
						},
					},
					CurrentPrice: &model.NotionPrice{
						Number: pointer.Ptr(uint64(7678)),
					},
					LowestPrice: &model.NotionPrice{
						Number: pointer.Ptr(uint64(7678)),
					},
					NotionReleaseDate: &model.NotionReleaseDate{
						NotionDate: &model.NotionDate{
							Start: "2024-11-15",
						},
					},
				},
			},
		}
		wantErr := errUnexpectedStatusCode
		if _, gotErr := wg.CreateNotionWishlistItem(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})
}

func TestUpdateNotionWishlistItem(t *testing.T) {
	t.Parallel()

	t.Run("Positive case: Successfully update a wishlist item in the Notion DB", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(_ *http.Request) (*http.Response, error) {
				jsonFile, err := os.Open("./testdata/updated_wishlist_item.json")
				if err != nil {
					t.Fatalf("failed to open updated_wishlist_item.json: %v", err)
				}
				defer jsonFile.Close()

				buffer := bytes.Buffer{}
				if _, err := io.Copy(&buffer, jsonFile); err != nil {
					t.Fatalf("failed to read updated_wishlist_item.json: %v", err)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(buffer.Bytes())),
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy_notion_api_key",
			NotionDatabaseID: "dummy_notion_database_id",
		}
		wg := NewNotionWishlistItemUpdater(cfg, m)
		input := &service.UpdateNotionWishlistItemInput{
			WishlistItem: &model.NotionWishlistItem{
				ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
				Parent: &model.NotionParent{
					DatabaseID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
				},
				Properties: &model.NotionProperties{
					NotionAppID: &model.NotionAppID{
						Title: []*model.NotionContent{
							{
								NotionText: &model.NotionText{
									NotionContent: "2701660",
								},
							},
						},
					},
					NotionName: &model.NotionName{
						RichText: []*model.NotionContent{
							{
								NotionText: &model.NotionText{
									NotionContent: "ドラゴンクエストIII　そして伝説へ…",
								},
							},
						},
					},
					CurrentPrice: &model.NotionPrice{
						Number: pointer.Ptr(uint64(7678)),
					},
					LowestPrice: &model.NotionPrice{
						Number: pointer.Ptr(uint64(7678)),
					},
					NotionReleaseDate: &model.NotionReleaseDate{
						NotionDate: &model.NotionDate{
							Start: "2024-11-15",
						},
					},
				},
			},
		}
		if _, err := wg.UpdateNotionWishlistItem(ctx, input); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Negative case: Fail to send a request", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		wantErr := errors.New("unexpected error")
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(nil, wantErr)

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy_notion_api_key",
			NotionDatabaseID: "dummy_notion_database_id",
		}
		wg := NewNotionWishlistItemUpdater(cfg, m)
		input := &service.UpdateNotionWishlistItemInput{
			WishlistItem: &model.NotionWishlistItem{
				ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
				Parent: &model.NotionParent{
					DatabaseID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
				},
				Properties: &model.NotionProperties{
					NotionAppID: &model.NotionAppID{
						Title: []*model.NotionContent{
							{
								NotionText: &model.NotionText{
									NotionContent: "2701660",
								},
							},
						},
					},
					NotionName: &model.NotionName{
						RichText: []*model.NotionContent{
							{
								NotionText: &model.NotionText{
									NotionContent: "ドラゴンクエストIII　そして伝説へ…",
								},
							},
						},
					},
					CurrentPrice: &model.NotionPrice{
						Number: pointer.Ptr(uint64(7678)),
					},
					LowestPrice: &model.NotionPrice{
						Number: pointer.Ptr(uint64(7678)),
					},
					NotionReleaseDate: &model.NotionReleaseDate{
						NotionDate: &model.NotionDate{
							Start: "2024-11-15",
						},
					},
				},
			},
		}
		if _, gotErr := wg.UpdateNotionWishlistItem(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Get a status code except 200", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       http.NoBody,
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy_notion_api_key",
			NotionDatabaseID: "dummy_notion_database_id",
		}
		wg := NewNotionWishlistItemUpdater(cfg, m)
		input := &service.UpdateNotionWishlistItemInput{
			WishlistItem: &model.NotionWishlistItem{
				ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
				Parent: &model.NotionParent{
					DatabaseID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
				},
				Properties: &model.NotionProperties{
					NotionAppID: &model.NotionAppID{
						Title: []*model.NotionContent{
							{
								NotionText: &model.NotionText{
									NotionContent: "2701660",
								},
							},
						},
					},
					NotionName: &model.NotionName{
						RichText: []*model.NotionContent{
							{
								NotionText: &model.NotionText{
									NotionContent: "ドラゴンクエストIII　そして伝説へ…",
								},
							},
						},
					},
					CurrentPrice: &model.NotionPrice{
						Number: pointer.Ptr(uint64(7678)),
					},
					LowestPrice: &model.NotionPrice{
						Number: pointer.Ptr(uint64(7678)),
					},
					NotionReleaseDate: &model.NotionReleaseDate{
						NotionDate: &model.NotionDate{
							Start: "2024-11-15",
						},
					},
				},
			},
		}
		wantErr := errUnexpectedStatusCode
		if _, gotErr := wg.UpdateNotionWishlistItem(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})
}

func TestDeleteNotionWishlistItem(t *testing.T) {
	t.Parallel()

	t.Run("Positive case: Successfully update a wishlist item in the Notion DB", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(_ *http.Request) (*http.Response, error) {
				jsonFile, err := os.Open("./testdata/updated_wishlist_item.json")
				if err != nil {
					t.Fatalf("failed to open updated_wishlist_item.json: %v", err)
				}
				defer jsonFile.Close()

				buffer := bytes.Buffer{}
				if _, err := io.Copy(&buffer, jsonFile); err != nil {
					t.Fatalf("failed to read updated_wishlist_item.json: %v", err)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(buffer.Bytes())),
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy_notion_api_key",
			NotionDatabaseID: "dummy_notion_database_id",
		}
		wg := NewNotionWishlistItemDeleter(cfg, m)
		input := &service.DeleteNotionWishlistItemInput{
			WishlistItem: &model.NotionWishlistItem{
				ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			},
		}
		if _, err := wg.DeleteNotionWishlistItem(ctx, input); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Negative case: Fail to send a request", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		wantErr := errors.New("unexpected error")
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(nil, wantErr)

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy_notion_api_key",
			NotionDatabaseID: "dummy_notion_database_id",
		}
		wg := NewNotionWishlistItemDeleter(cfg, m)
		input := &service.DeleteNotionWishlistItemInput{
			WishlistItem: &model.NotionWishlistItem{
				ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			},
		}
		if _, gotErr := wg.DeleteNotionWishlistItem(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Get a status code except 200", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       http.NoBody,
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.NotionConfig{
			NotionAPIKey:     "dummy_notion_api_key",
			NotionDatabaseID: "dummy_notion_database_id",
		}
		wg := NewNotionWishlistItemDeleter(cfg, m)
		input := &service.DeleteNotionWishlistItemInput{
			WishlistItem: &model.NotionWishlistItem{
				ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			},
		}
		wantErr := errUnexpectedStatusCode
		if _, gotErr := wg.DeleteNotionWishlistItem(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})
}
