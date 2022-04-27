package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func initLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func setLogLevel(logLevelStr string) {
	logLevel, err := zerolog.ParseLevel(logLevelStr)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	zerolog.SetGlobalLevel(logLevel)
}
