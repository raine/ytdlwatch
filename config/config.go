package config

import (
	"os"

	"github.com/rs/zerolog/log"
)

type Config struct {
	LogLevel      string
	OutputPath    string
	YoutubedlPath string

	YoutubeApiKey     string
	YoutubePlaylistId string

	PlexApiToken   string
	PlexApiUrl     string
	PlexLibraryKey string
}

func GetConfig() Config {
	logLevel, _ := os.LookupEnv("LOG_LEVEL")
	outputPath, _ := os.LookupEnv("OUTPUT_PATH")
	youtubedlBinPath, _ := os.LookupEnv("YOUTUBE_DL_PATH")

	youtubeApiKey, _ := os.LookupEnv("YOUTUBE_API_KEY")
	youtubePlaylistId, _ := os.LookupEnv("YOUTUBE_PLAYLIST_ID")

	plexApiToken, _ := os.LookupEnv("PLEX_API_TOKEN")
	plexApiUrl, _ := os.LookupEnv("PLEX_API_URL")
	plexLibraryKey, _ := os.LookupEnv("PLEX_LIBRARY_KEY")

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
		LogLevel:      logLevel,
		OutputPath:    outputPath,
		YoutubedlPath: youtubedlBinPath,

		YoutubeApiKey:     youtubeApiKey,
		YoutubePlaylistId: youtubePlaylistId,

		PlexApiToken:   plexApiToken,
		PlexApiUrl:     plexApiUrl,
		PlexLibraryKey: plexLibraryKey,
	}
}
