package plex

import (
	"strings"
	"time"

	"github.com/raine/go-plex-client"
	"github.com/rs/zerolog/log"
)

type PlexWatchedVidDeleter struct {
	plexClient      *plex.Plex
	plexLibraryKey  string
	plexLibraryPath string
	ticker          *time.Ticker
}

func (p *PlexWatchedVidDeleter) getWatchedVideosFromLibrary() []plex.Metadata {
	results, err := p.plexClient.GetLibraryContent(p.plexLibraryKey, "")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get plex library content")
	}

	watched := make([]plex.Metadata, 0)

	for _, metadata := range results.MediaContainer.Metadata {
		viewCount := metadata.ViewCount.String()
		if viewCount != "" {
			watched = append(watched, metadata)
		}
	}

	return watched
}

func (p *PlexWatchedVidDeleter) deleteVideo(video plex.Metadata) {
	if len(video.Media) == 0 {
		log.Error().Msg("expected plex video to have media")
		return
	}
	media := video.Media[0]

	if len(media.Part) == 0 {
		log.Error().Msg("expected plex media to have parts")
		return
	}
	part := media.Part[0].File

	if !strings.HasPrefix(part, p.plexLibraryPath) {
		log.Error().
			Str("plexLibraryPath", p.plexLibraryPath).
			Msg("media part's path does not start with configured youtube-dl output path, not removing in case the plex library key is incorrectly configured")
		return
	}

	err := p.plexClient.DeleteMediaByID(video.RatingKey)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete video")
		return
	}

	log.Info().Str("title", video.Title).Msg("deleted video")
}

func (p *PlexWatchedVidDeleter) deleteWatchedVideos() {
	watchedVideos := p.getWatchedVideosFromLibrary()
	if len(watchedVideos) > 0 {
		log.Info().Msgf("will delete %d watched videos", len(watchedVideos))

		for _, v := range watchedVideos {
			p.deleteVideo(v)
		}
	}
}

func (p *PlexWatchedVidDeleter) Start() {
	p.ticker = time.NewTicker(time.Minute * 10)

	go func() {
		for range p.ticker.C {
			p.deleteWatchedVideos()
		}
	}()
}

func NewPlexWatchedVidDeleter(
	plexApiUrl string,
	plexApiToken string,
	plexLibraryKey string,
	plexLibraryPath string,
) *PlexWatchedVidDeleter {
	plexClient, err := plex.New(plexApiUrl, plexApiToken)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create plex client")
	}

	_, err = plexClient.Test()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to the plex server")
	}

	log.Info().Msg("connected to plex api successfully")

	p := PlexWatchedVidDeleter{
		plexClient:      plexClient,
		plexLibraryKey:  plexLibraryKey,
		plexLibraryPath: plexLibraryPath,
	}

	return &p
}
