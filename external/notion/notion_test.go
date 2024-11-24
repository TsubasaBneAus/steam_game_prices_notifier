package notion

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/TsubasaBneAus/steam_game_price_notifier/config"
	httpClient "github.com/TsubasaBneAus/steam_game_price_notifier/external/httpClient/mock"
	"github.com/TsubasaBneAus/steam_game_price_notifier/model"
	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/pointer"
	"go.uber.org/mock/gomock"
)

func TestGetWishlist(t *testing.T) {
	t.Parallel()

	t.Run("Positive case: Get a wishlist successfully", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpClient.NewMockHttpClient(ctrl)
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
		wg := NewWishlistGetter(&config.Envs{}, m)
		got, err := wg.GetWishlist(ctx, &GetWishlistInput{})
		if err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
		want := &GetWishlistOutput{
			WishlistItems: &model.WishlistItems{
				Results: []*model.WishlistItem{
					{
						ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
						Properties: &model.Properties{
							AppID: &model.AppID{
								Title: []*model.Content{},
							},
							Name: &model.Name{
								RichText: []*model.Content{},
							},
							CurrentPrice: &model.Price{
								Number: nil,
							},
							LowestPrice: &model.Price{
								Number: nil,
							},
							ReleaseDate: &model.ReleaseDate{
								Date: nil,
							},
						},
					},
					{
						ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
						Properties: &model.Properties{
							AppID: &model.AppID{
								Title: []*model.Content{
									{
										Text: &model.Text{
											Content: "2701660",
										},
									},
								},
							},
							Name: &model.Name{
								RichText: []*model.Content{
									{
										Text: &model.Text{
											Content: "ドラゴンクエストIII　そして伝説へ…",
										},
									},
								},
							},
							CurrentPrice: &model.Price{
								Number: pointer.Ptr(uint64(7678)),
							},
							LowestPrice: &model.Price{
								Number: pointer.Ptr(uint64(7678)),
							},
							ReleaseDate: &model.ReleaseDate{
								Date: &model.Date{
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

	t.Run("Positive case: Get an empty wishlist successfully", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpClient.NewMockHttpClient(ctrl)
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
		wg := NewWishlistGetter(&config.Envs{}, m)
		got, err := wg.GetWishlist(ctx, &GetWishlistInput{})
		if err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
		want := &GetWishlistOutput{
			WishlistItems: &model.WishlistItems{
				Results: []*model.WishlistItem{},
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
		m := httpClient.NewMockHttpClient(ctrl)
		wantErr := errors.New("unexpected error")
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(nil, wantErr)

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		wg := NewWishlistGetter(&config.Envs{}, m)
		if _, gotErr := wg.GetWishlist(ctx, &GetWishlistInput{}); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Get a status code except 200", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpClient.NewMockHttpClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		wg := NewWishlistGetter(&config.Envs{}, m)
		wantErr := errUnexpectedStatusCode
		if _, gotErr := wg.GetWishlist(ctx, &GetWishlistInput{}); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Fail to unmarshal a response", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpClient.NewMockHttpClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("invalid JSON"))),
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		wg := NewWishlistGetter(&config.Envs{}, m)
		if _, err := wg.GetWishlist(ctx, &GetWishlistInput{}); err == nil {
			t.Errorf("\ngot: %v\nwant: an error generated by the library", nil)
		}
	})
}

func TestUpdateWishlist(t *testing.T) {
	t.Parallel()

	t.Run("Positive case: Update a wishlist successfully", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpClient.NewMockHttpClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		wg := NewWishlistUpdater(&config.Envs{}, m)
		input := &UpdateWishlistInput{
			WishlistItem: &model.WishlistItem{
				ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
				Properties: &model.Properties{
					AppID: &model.AppID{
						Title: []*model.Content{
							{
								Text: &model.Text{
									Content: "2701660",
								},
							},
						},
					},
					Name: &model.Name{
						RichText: []*model.Content{
							{
								Text: &model.Text{
									Content: "ドラゴンクエストIII　そして伝説へ…",
								},
							},
						},
					},
					CurrentPrice: &model.Price{
						Number: pointer.Ptr(uint64(7678)),
					},
					LowestPrice: &model.Price{
						Number: pointer.Ptr(uint64(7678)),
					},
					ReleaseDate: &model.ReleaseDate{
						Date: &model.Date{
							Start: "2024-11-15",
						},
					},
				},
			},
		}
		if _, err := wg.UpdateWishlist(ctx, input); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Positive case: Update a wishlist with an empty input successfully", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpClient.NewMockHttpClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		wg := NewWishlistUpdater(&config.Envs{}, m)
		input := &UpdateWishlistInput{
			WishlistItem: &model.WishlistItem{
				ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
				Properties: &model.Properties{
					AppID:        nil,
					Name:         nil,
					CurrentPrice: nil,
					LowestPrice:  nil,
					ReleaseDate:  nil,
				},
			},
		}
		if _, err := wg.UpdateWishlist(ctx, input); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
	})

	t.Run("Negative case: Fail to send a request", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpClient.NewMockHttpClient(ctrl)
		wantErr := errors.New("unexpected error")
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(nil, wantErr)

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		wg := NewWishlistUpdater(&config.Envs{}, m)
		input := &UpdateWishlistInput{
			WishlistItem: &model.WishlistItem{
				ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
				Properties: &model.Properties{
					AppID: &model.AppID{
						Title: []*model.Content{
							{
								Text: &model.Text{
									Content: "2701660",
								},
							},
						},
					},
					Name: &model.Name{
						RichText: []*model.Content{
							{
								Text: &model.Text{
									Content: "ドラゴンクエストIII　そして伝説へ…",
								},
							},
						},
					},
					CurrentPrice: &model.Price{
						Number: pointer.Ptr(uint64(7678)),
					},
					LowestPrice: &model.Price{
						Number: pointer.Ptr(uint64(7678)),
					},
					ReleaseDate: &model.ReleaseDate{
						Date: &model.Date{
							Start: "2024-11-15",
						},
					},
				},
			},
		}
		if _, gotErr := wg.UpdateWishlist(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Get a status code except 200", func(t *testing.T) {
		t.Parallel()

		// Create a mock of the HTTP client
		ctrl := gomock.NewController(t)
		m := httpClient.NewMockHttpClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		wg := NewWishlistUpdater(&config.Envs{}, m)
		input := &UpdateWishlistInput{
			WishlistItem: &model.WishlistItem{
				ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
				Properties: &model.Properties{
					AppID: &model.AppID{
						Title: []*model.Content{
							{
								Text: &model.Text{
									Content: "2701660",
								},
							},
						},
					},
					Name: &model.Name{
						RichText: []*model.Content{
							{
								Text: &model.Text{
									Content: "ドラゴンクエストIII　そして伝説へ…",
								},
							},
						},
					},
					CurrentPrice: &model.Price{
						Number: pointer.Ptr(uint64(7678)),
					},
					LowestPrice: &model.Price{
						Number: pointer.Ptr(uint64(7678)),
					},
					ReleaseDate: &model.ReleaseDate{
						Date: &model.Date{
							Start: "2024-11-15",
						},
					},
				},
			},
		}
		wantErr := errUnexpectedStatusCode
		if _, gotErr := wg.UpdateWishlist(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})
}
