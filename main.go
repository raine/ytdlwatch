package main

import (
	"os"
	"runtime"

	"github.com/raine/ytdlwatch/config"
	"github.com/raine/ytdlwatch/plex"
	"github.com/raine/ytdlwatch/youtube"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var sem = make(chan int, 1)

func processVideoUrls(config config.Config, videoUrls chan string) {
	for url := range videoUrls {
		url := url
		log.Info().Str("url", url).Msg("got a video url")
		sem <- 1
		go func() {
			download(
				config.YoutubedlPath,
				config.YoutubedlArgs,
				url,
			)
			<-sem
		}()
	}
}

func main() {
	handleGracefulExit()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	videoUrls := make(chan string)
	config := config.GetConfig()

	logLevel, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	zerolog.SetGlobalLevel(logLevel)

	go processVideoUrls(config, videoUrls)

	if config.YoutubePlaylistId != "" {
		log.Info().Str("youtubePlaylistId", config.YoutubePlaylistId).Msg("youtube playlist configured")

		if config.YoutubeApiKey != "" {
			poller := youtube.NewYoutubePlaylistPoller(config.YoutubeApiKey, config.YoutubePlaylistId, videoUrls)
			go poller.Start()
		} else {
			log.Fatal().Msg("expected YOUTUBE_API_KEY environment variable to be set with YOUTUBE_PLAYLIST_ID")
		}
	} else {
		log.Fatal().Msg("youtube playlist not configured")
	}

	if config.PlexApiToken != "" &&
		config.PlexApiUrl != "" &&
		config.PlexLibraryKey != "" {
		log.Info().
			Str("plexApiUrl", config.PlexApiUrl).
			Str("plexLibraryKey", config.PlexLibraryKey).
			Msg("plex configured, will delete watched videos")

		plexVideoDeleter := plex.NewPlexWatchedVidDeleter(
			config.PlexApiUrl,
			config.PlexApiToken,
			config.PlexLibraryKey,
			config.OutputPath,
		)
		plexVideoDeleter.Start()
	} else {
		log.Info().Msg("plex not configured")
	}

	runtime.Goexit()
}
