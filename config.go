package main

import (
	"os"

	"github.com/rs/zerolog/log"
)

type Config struct {
	youtubeApiKey     string
	youtubePlaylistId string
	outputPath        string
	youtubedlPath     string
	logLevel          string
}

func getConfig() Config {
	youtubeApiKey, _ := os.LookupEnv("YOUTUBE_API_KEY")
	youtubePlaylistId, _ := os.LookupEnv("YOUTUBE_PLAYLIST_ID")
	outputPath, _ := os.LookupEnv("OUTPUT_PATH")
	youtubedlBinPath, _ := os.LookupEnv("YOUTUBE_DL_PATH")
	logLevel, _ := os.LookupEnv("LOG_LEVEL")

	if outputPath == "" {
		log.Fatal().Msg("expected OUTPUT_PATH to be set")
	}

	if youtubedlBinPath == "" {
		youtubedlBinPath = "yt-dlp"
		log.Info().Msgf("YOUTUBE_DL_PATH is unset, defaulting to %s", youtubedlBinPath)
	}

	if logLevel == "" {
		logLevel = "info"
	}

	return Config{
		youtubeApiKey:     youtubeApiKey,
		youtubePlaylistId: youtubePlaylistId,
		outputPath:        outputPath,
		youtubedlPath:     youtubedlBinPath,
		logLevel:          logLevel,
	}
}
