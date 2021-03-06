package main

import (
	"bufio"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func download(
	youtubedlPath string,
	youtubedlArgs []string,
	videoUrl string,
) {
	args := append(youtubedlArgs, videoUrl)

	cmd := exec.Command(
		youtubedlPath,
		args...,
	)
	log.Info().Msgf("running %+v", cmd)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	go func() {
		for stdoutScanner.Scan() {
			log.Info().Msgf("yt-dlp: %s", stdoutScanner.Text())
		}
	}()

	go func() {
		for stderrScanner.Scan() {
			log.Info().Msgf("yt-dlp: %s", stderrScanner.Text())
		}
	}()

	if err := cmd.Start(); err != nil {
		log.Error().Err(err).Send()
		return
	}

	if err := cmd.Wait(); err != nil {
		log.Error().Err(err).Msg("failed to run youtube-dl")
	}

	log.Info().Msgf("finished running %+v", cmd)
}
