package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/TsubasaBneAus/steam_game_price_notifier/app/usecase"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context) {
	// Set a JSON formatted logger as the default logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Initialize the application
	app, err := InitializeApp(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to initialize the application", slog.Any("error", err))
		os.Exit(1)
	}

	// Notify video game prices
	if _, err := app.vGPNotifier.NotifyVideoGamePrices(ctx, &usecase.NotifyVideoGamePricesInput{}); err != nil {
		slog.ErrorContext(ctx, "failed to notify video game prices", slog.Any("error", err))

		// Notify an error
		input := &usecase.NotifyErrorInput{
			GeneratedError: err,
		}
		if _, err := app.eNotifier.NotifyError(ctx, input); err != nil {
			slog.ErrorContext(ctx, "failed to notify an error", slog.Any("error", err))
		}

		os.Exit(1)
	}
}
