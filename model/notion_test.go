package model

import (
	"context"
	"testing"
)

func TestNotionToTime(t *testing.T) {
	t.Parallel()

	t.Run("Positive case: Convert date string successfully", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		date := Date{
			Start: "2024-12-31",
		}
		if _, err := date.ToTime(ctx); err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
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

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			date := Date{
				Start: tc.date,
			}
			if _, err := date.ToTime(ctx); err == nil {
				t.Errorf("\ngot: %v\nwant: %v", err, nil)
			}
		})
	}
}
