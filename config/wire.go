package config

import "github.com/google/wire"

// A wire set for the config package
var Set = wire.NewSet(
	NewNotionConfig,
	NewSteamConfig,
	NewDiscordConfig,
)
