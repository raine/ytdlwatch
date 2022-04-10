package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func handleGracefulExit() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM)

	go func() {
		s := <-sigc
		log.Info().Msgf("got %s, exiting", s)
		os.Exit(1)
	}()
}
