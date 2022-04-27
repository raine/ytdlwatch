package config

import (
	"os"

	"github.com/mattn/go-shellwords"
	"github.com/rs/zerolog/log"
)

const (
	defaultLogLevel         = "info"
	defaultYoutubedlBinPath = "yt-dlp"
)

type Config struct {
	LogLevel          string
	OutputPath        string
	YoutubedlPath     string
	YoutubedlArgs     []string
	YoutubeApiKey     string
	YoutubePlaylistId string
	PlexApiToken      string
	PlexApiUrl        string
	PlexLibraryKey    string
}

func GetConfig() Config {
	logLevel, _ := os.LookupEnv("LOG_LEVEL")
	outputPath, _ := os.LookupEnv("OUTPUT_PATH")
	youtubedlBinPath, _ := os.LookupEnv("YOUTUBE_DL_PATH")
	youtubedlArgsRaw, _ := os.LookupEnv("YOUTUBE_DL_ARGS")

	youtubeApiKey, _ := os.LookupEnv("YOUTUBE_API_KEY")
	youtubePlaylistId, _ := os.LookupEnv("YOUTUBE_PLAYLIST_ID")

	plexApiToken, _ := os.LookupEnv("PLEX_API_TOKEN")
	plexApiUrl, _ := os.LookupEnv("PLEX_API_URL")
	plexLibraryKey, _ := os.LookupEnv("PLEX_LIBRARY_KEY")

	if outputPath == "" {
		log.Fatal().Msg("expected OUTPUT_PATH to be set")
	}

	if youtubedlBinPath == "" {
		youtubedlBinPath = defaultYoutubedlBinPath
		log.Info().Msgf("YOUTUBE_DL_PATH is unset, defaulting to %s", youtubedlBinPath)
	}

	var youtubedlArgs []string
	if youtubedlArgsRaw == "" {
		youtubedlArgs = []string{
			"--newline",
			"--output", "%(title)s.%(ext)s",
			"--embed-metadata",
			"--embed-subs",
			"--sub-lang", "en",
			"--paths", outputPath,
		}
	} else {
		args, err := shellwords.Parse(youtubedlArgsRaw)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to parse arguments from YOUTUBE_DL_ARGS env var")
		}
		youtubedlArgs = args
	}

	if logLevel == "" {
		logLevel = defaultLogLevel
	}

	return Config{
		LogLevel:      logLevel,
		OutputPath:    outputPath,
		YoutubedlPath: youtubedlBinPath,
		YoutubedlArgs: youtubedlArgs,

		YoutubeApiKey:     youtubeApiKey,
		YoutubePlaylistId: youtubePlaylistId,

		PlexApiToken:   plexApiToken,
		PlexApiUrl:     plexApiUrl,
		PlexLibraryKey: plexLibraryKey,
	}
}
