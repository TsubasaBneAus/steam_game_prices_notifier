package usecase

import (
	"context"
)

type (
	// An input to notify video game prices on Discord
	NotifyVideoGamePricesInput struct{}

	// An output to notify video game prices on Discord
	NotifyVideoGamePricesOutput struct{}

	// An interface to notify video game prices on Discord
	VideoGamePricesNotifier interface {
		NotifyVideoGamePrices(
			ctx context.Context,
			input *NotifyVideoGamePricesInput,
		) (*NotifyVideoGamePricesOutput, error)
	}
)

type (
	// An input to notify an error on Discord
	NotifyErrorInput struct {
		GeneratedError error
	}

	// An output to notify an error on Discord
	NotifyErrorOutput struct{}

	// An interface to notify an error on Discord
	ErrorNotifier interface {
		NotifyError(ctx context.Context, input *NotifyErrorInput) (*NotifyErrorOutput, error)
	}
)
