package main

import "github.com/TsubasaBneAus/steam_game_price_notifier/app/usecase"

type app struct {
	vGPNotifier usecase.VideoGamePricesNotifier
	eNotifier   usecase.ErrorNotifier
}

// Generate a new app
func NewApp(
	vGPNotifier usecase.VideoGamePricesNotifier,
	eNotifier usecase.ErrorNotifier,
) *app {
	return &app{
		vGPNotifier: vGPNotifier,
		eNotifier:   eNotifier,
	}
}
