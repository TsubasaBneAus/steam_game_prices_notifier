package model

import (
	"context"
	"testing"
)

func TestToUint64(t *testing.T) {
	t.Parallel()

	t.Run("Positive case: Successfully convert json.Number into uint64", func(t *testing.T) {
		t.Parallel()

		// Execute the method to be tested
		ctx := t.Context()
		currentPrice := SteamCurrentPrice{Number: "100000"}
		got, err := currentPrice.ConvertPriceFormat(ctx)
		if err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
		want := uint64(1000)
		if *got != want {
			t.Errorf("\ngot: %v\nwant: %v", *got, want)
		}
	})

	t.Run("Negative case: Failed to convert json.Number into int64", func(t *testing.T) {
		t.Parallel()

		// Execute the method to be tested
		ctx := t.Context()
		currentPrice := SteamCurrentPrice{Number: "9223372036854775808"}
		if _, err := currentPrice.ConvertPriceFormat(ctx); err == nil {
			t.Errorf("\ngot: %v\nwant: an error generated by the library", nil)
		}
	})
}

func TestToTime(t *testing.T) {
	t.Parallel()

	t.Run("Positive case: Successfully convert json.Number into time.Time", func(t *testing.T) {
		t.Parallel()

		// Execute the method to be tested
		ctx := t.Context()
		releaseDate := SteamReleaseDate{
			Date: "11 Nov, 2021",
		}
		got, err := releaseDate.ToTime(ctx)
		if err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
		}
		want := "2021-11-11 00:00:00 +0900 JST"
		if got.String() != want {
			t.Errorf("\ngot: %v\nwant: %v", got, want)
		}
	})

	negativeTestCases := map[string]struct {
		date string
	}{
		"Negative case: The date string to be tested is invalid": {
			date: "2024-13-32",
		},
		"Negative case: The date string to be tested is empty": {
			date: "",
		},
	}

	for name, tc := range negativeTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Execute the method to be tested
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			releaseDate := SteamReleaseDate{
				Date: tc.date,
			}
			if _, err := releaseDate.ToTime(ctx); err == nil {
				t.Errorf("\ngot: %v\nwant: %v", err, nil)
			}
		})
	}
}
