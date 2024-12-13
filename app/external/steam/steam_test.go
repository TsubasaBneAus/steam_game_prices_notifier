package steam

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
	"go.uber.org/mock/gomock"
)

func TestGetSteamWishlist(t *testing.T) {
	t.Parallel()

	t.Run("Positive case: Successfully get a wishlist from the Steam Store", func(t *testing.T) {
		t.Parallel()

		// Create a mock for the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(req *http.Request) (*http.Response, error) {
				got := req.URL.String()
				want := "https://api.steampowered.com/IWishlistService/GetWishlist/v1/?steamid=dummy_steam_user_id"
				if diff := cmp.Diff(got, want); diff != "" {
					t.Errorf("got(-) want(+)\n%s", diff)
				}

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
		cfg := &config.SteamConfig{
			SteamUserID: "dummy_steam_user_id",
		}
		wg := NewSteamWishlistGetter(cfg, m)
		got, err := wg.GetSteamWishlist(ctx, &service.GetSteamWishlistInput{})
		if err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
		want := &service.GetSteamWishlistOutput{
			Wishlist: &model.SteamStoreWishlist{
				Response: &model.SteamStoreResponse{
					Items: []*model.SteamStoreItem{
						{
							AppID: 105600,
						},
						{
							AppID: 251570,
						},
					},
				},
			},
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("got(-) want(+)\n%s", diff)
		}
	})

	t.Run("Positive case: Successfully get an empty wishlist from the Steam Store", func(t *testing.T) {
		t.Parallel()

		// Create a mock for the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(req *http.Request) (*http.Response, error) {
				got := req.URL.String()
				want := "https://api.steampowered.com/IWishlistService/GetWishlist/v1/?steamid=dummy_steam_user_id"
				if diff := cmp.Diff(got, want); diff != "" {
					t.Errorf("got(-) want(+)\n%s", diff)
				}

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
		cfg := &config.SteamConfig{
			SteamUserID: "dummy_steam_user_id",
		}
		wg := NewSteamWishlistGetter(cfg, m)
		got, err := wg.GetSteamWishlist(ctx, &service.GetSteamWishlistInput{})
		if err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
		want := &service.GetSteamWishlistOutput{
			Wishlist: &model.SteamStoreWishlist{
				Response: &model.SteamStoreResponse{
					Items: []*model.SteamStoreItem{},
				},
			},
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("got(-) want(+)\n%s", diff)
		}
	})

	t.Run("Negative case: Fail to send a request", func(t *testing.T) {
		t.Parallel()

		// Create a mock for the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		wantErr := errors.New("unexpected error")
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(req *http.Request) (*http.Response, error) {
				got := req.URL.String()
				want := "https://api.steampowered.com/IWishlistService/GetWishlist/v1/?steamid=dummy_steam_user_id"
				if diff := cmp.Diff(got, want); diff != "" {
					t.Errorf("got(-) want(+)\n%s", diff)
				}

				return nil, wantErr
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.SteamConfig{
			SteamUserID: "dummy_steam_user_id",
		}
		wg := NewSteamWishlistGetter(cfg, m)
		if _, gotErr := wg.GetSteamWishlist(ctx, &service.GetSteamWishlistInput{}); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Get a status code except 200", func(t *testing.T) {
		t.Parallel()

		// Create a mock for the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(req *http.Request) (*http.Response, error) {
				got := req.URL.String()
				want := "https://api.steampowered.com/IWishlistService/GetWishlist/v1/?steamid=dummy_steam_user_id"
				if diff := cmp.Diff(got, want); diff != "" {
					t.Errorf("got(-) want(+)\n%s", diff)
				}

				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       http.NoBody,
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.SteamConfig{
			SteamUserID: "dummy_steam_user_id",
		}
		wg := NewSteamWishlistGetter(cfg, m)
		wantErr := errUnexpectedStatusCode
		if _, gotErr := wg.GetSteamWishlist(ctx, &service.GetSteamWishlistInput{}); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Fail to unmarshal a response", func(t *testing.T) {
		t.Parallel()

		// Create a mock for the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(req *http.Request) (*http.Response, error) {
				got := req.URL.String()
				want := "https://api.steampowered.com/IWishlistService/GetWishlist/v1/?steamid=dummy_steam_user_id"
				if diff := cmp.Diff(got, want); diff != "" {
					t.Errorf("got(-) want(+)\n%s", diff)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("invalid JSON"))),
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.SteamConfig{
			SteamUserID: "dummy_steam_user_id",
		}
		wg := NewSteamWishlistGetter(cfg, m)
		if _, err := wg.GetSteamWishlist(ctx, &service.GetSteamWishlistInput{}); err == nil {
			t.Errorf("\ngot: %v\nwant: an error generated by the library", nil)
		}
	})
}

func TestGetSteamVideoGameDetails(t *testing.T) {
	t.Parallel()

	t.Run("Positive case: Successfully get video game details from the Steam Store", func(t *testing.T) {
		t.Parallel()

		// Create a mock for the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(req *http.Request) (*http.Response, error) {
				got := req.URL.String()
				want := "https://store.steampowered.com/api/appdetails/?appids=2701660&cc=jp"
				if diff := cmp.Diff(got, want); diff != "" {
					t.Errorf("got(-) want(+)\n%s", diff)
				}

				jsonFile, err := os.Open("./testdata/video_game_details.json")
				if err != nil {
					t.Fatalf("failed to open video_game_details.json: %v", err)
				}
				defer jsonFile.Close()

				buffer := bytes.Buffer{}
				if _, err := io.Copy(&buffer, jsonFile); err != nil {
					t.Fatalf("failed to read video_game_details.json: %v", err)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(buffer.Bytes())),
				}, nil
			})

		// Execute the method to be tested (Skip checking the response)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.SteamConfig{
			SteamUserID: "dummy_steam_user_id",
		}
		vg := NewSteamVideoGameDetailsGetter(cfg, m)
		input := &service.GetSteamVideoGameDetailsInput{
			AppID: 2701660,
		}
		got, err := vg.GetSteamVideoGameDetails(ctx, input)
		if err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
		want := &service.GetSteamVideoGameDetailsOutput{
			VideoGameDetails: &model.SteamStoreVideoGameDetails{
				AppID: 2701660,
				Title: "DRAGON QUEST III HD-2D Remake",
				CurrentPrice: &model.SteamCurrentPrice{
					Number: "767800",
				},
				ReleaseDate: &model.SteamReleaseDate{
					Date: "14 Nov, 2024",
				},
			},
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("got(-) want(+)\n%s", diff)
		}
	})

	t.Run("Negative case: Fail to send a request", func(t *testing.T) {
		t.Parallel()

		// Create a mock for the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		wantErr := errors.New("unexpected error")
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(req *http.Request) (*http.Response, error) {
				got := req.URL.String()
				want := "https://store.steampowered.com/api/appdetails/?appids=2701660&cc=jp"
				if diff := cmp.Diff(got, want); diff != "" {
					t.Errorf("got(-) want(+)\n%s", diff)
				}

				return nil, wantErr
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.SteamConfig{
			SteamUserID: "dummy_steam_user_id",
		}
		vg := NewSteamVideoGameDetailsGetter(cfg, m)
		input := &service.GetSteamVideoGameDetailsInput{
			AppID: 2701660,
		}
		if _, gotErr := vg.GetSteamVideoGameDetails(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Get a status code except 200", func(t *testing.T) {
		t.Parallel()

		// Create a mock for the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(req *http.Request) (*http.Response, error) {
				got := req.URL.String()
				want := "https://store.steampowered.com/api/appdetails/?appids=2701660&cc=jp"
				if diff := cmp.Diff(got, want); diff != "" {
					t.Errorf("got(-) want(+)\n%s", diff)
				}

				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       http.NoBody,
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.SteamConfig{
			SteamUserID: "dummy_steam_user_id",
		}
		vg := NewSteamVideoGameDetailsGetter(cfg, m)
		input := &service.GetSteamVideoGameDetailsInput{
			AppID: 2701660,
		}
		wantErr := errUnexpectedStatusCode
		if _, gotErr := vg.GetSteamVideoGameDetails(ctx, input); !errors.Is(gotErr, wantErr) {
			t.Errorf("\ngot: %v\nwant: %v", gotErr, wantErr)
		}
	})

	t.Run("Negative case: Fail to unmarshal a response", func(t *testing.T) {
		t.Parallel()

		// Create a mock for the HTTP client
		ctrl := gomock.NewController(t)
		m := httpclient.NewMockHTTPClient(ctrl)
		m.
			EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(req *http.Request) (*http.Response, error) {
				got := req.URL.String()
				want := "https://store.steampowered.com/api/appdetails/?appids=2701660&cc=jp"
				if diff := cmp.Diff(got, want); diff != "" {
					t.Errorf("got(-) want(+)\n%s", diff)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("invalid JSON"))),
				}, nil
			})

		// Execute the method to be tested
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg := &config.SteamConfig{
			SteamUserID: "dummy_steam_user_id",
		}
		vg := NewSteamVideoGameDetailsGetter(cfg, m)
		input := &service.GetSteamVideoGameDetailsInput{
			AppID: 2701660,
		}
		if _, err := vg.GetSteamVideoGameDetails(ctx, input); err == nil {
			t.Errorf("\ngot: %v\nwant: an error generated by the library", nil)
		}
	})
}
