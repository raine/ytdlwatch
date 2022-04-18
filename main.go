package main

import (
	"os"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var sem = make(chan int, 1)

func process(config Config, videoUrls chan string) {
	for url := range videoUrls {
		url := url
		log.Info().Str("url", url).Msg("got a video url")
		sem <- 1
		go func() {
			download(config, url)
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
	config := getConfig()

	logLevel, err := zerolog.ParseLevel(config.logLevel)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	zerolog.SetGlobalLevel(logLevel)

	go process(config, videoUrls)

	if config.youtubePlaylistId != "" {
		log.Info().Str("youtubePlaylistId", config.youtubePlaylistId).Msg("youtube playlist configured")

		if config.youtubeApiKey != "" {
			poller := NewYoutubePlaylistPoller(config.youtubeApiKey, config.youtubePlaylistId, videoUrls)
			poller.Start()
		} else {
			log.Fatal().Msg("expected YOUTUBE_API_KEY environment variable to be set with YOUTUBE_PLAYLIST_ID")
		}
	} else {
		log.Fatal().Msg("youtube playlist not configured")
	}

	runtime.Goexit()
}
