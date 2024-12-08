package usecase

import (
	"context"
)

// An input to notify video game prices on Discord
type NotifyVideoGamePricesInput struct{}

// An output to notify video game prices on Discord
type NotifyVideoGamePricesOutput struct{}

// An interface to notify video game prices on Discord
type VideoGamePricesNotifier interface {
	NotifyVideoGamePrices(
		ctx context.Context,
		input *NotifyVideoGamePricesInput,
	) (*NotifyVideoGamePricesOutput, error)
}
