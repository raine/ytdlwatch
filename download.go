package main

import (
	"bufio"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func download(
	config Config,
	videoUrl string,
) {
	cmd := exec.Command(
		config.youtubedlPath,
		"--newline",
		"--output", "%(title)s.%(ext)s",
		"--embed-metadata",
		"--embed-subs",
		"--sub-lang", "en",
		"--paths", config.outputPath,
		videoUrl,
	)
	log.Info().Msgf("running %+v", cmd)

	r, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	cmd.Stderr = cmd.Stdout
	scanner := bufio.NewScanner(r)

	go func() {
		for scanner.Scan() {
			log.Debug().Msgf("yt-dlp: %s", scanner.Text())
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
