package youtube

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubePlaylistPoller struct {
	youtubeApiKey     string
	youtubePlaylistId string
	videoUrls         chan string
	ticker            *time.Ticker
	prevCheckIdMap    map[string]struct{}
}

func getPlaylistItemVideoUrl(item *youtube.PlaylistItem) string {
	return fmt.Sprintf(
		"https://www.youtube.com/watch?v=%s",
		item.Snippet.ResourceId.VideoId,
	)
}

func (p *YoutubePlaylistPoller) getPlaylistItemList(maxResults int64) (*youtube.PlaylistItemListResponse, error) {
	log.Debug().Msg("getting playlist items")

	ctx := context.Background()
	youtube, err := youtube.NewService(ctx, option.WithAPIKey(p.youtubeApiKey))
	if err != nil {
		return nil, err
	}

	return youtube.PlaylistItems.
		List([]string{"snippet"}).
		MaxResults(maxResults).
		PlaylistId(p.youtubePlaylistId).
		Do()
}

func (p *YoutubePlaylistPoller) getNewPlaylistItems(maxResults int64) ([]*youtube.PlaylistItem, error) {
	list, err := p.getPlaylistItemList(maxResults)
	if err != nil {
		return nil, err
	}

	newItems := make([]*youtube.PlaylistItem, 0)

	for _, item := range list.Items {
		if _, ok := p.prevCheckIdMap[item.Id]; !ok {
			newItems = append(newItems, item)
		}

		// Collect all previously seen ids for now, not just ones from last query,
		// which has the problem that deleting videos from the playlist might make
		// old videos get returned
		p.prevCheckIdMap[item.Id] = struct{}{}
	}

	return newItems, nil
}

func (p *YoutubePlaylistPoller) Start() {
	log.Info().Msg("starting to poll a youtube playlist")
	p.getNewPlaylistItems(50)
	p.ticker = time.NewTicker(time.Second * 10)

	for range p.ticker.C {
		items, err := p.getNewPlaylistItems(25)
		if err != nil {
			log.Error().Err(err).Msg("failed to get new playlist items")
			continue
		}

		if len(items) > 0 {
			log.Info().Msgf("got %d new playlist items", len(items))
		}

		for _, item := range items {
			p.videoUrls <- getPlaylistItemVideoUrl(item)
		}
	}
}

func NewYoutubePlaylistPoller(
	youtubeApiKey string,
	youtubePlaylistId string,
	videoUrls chan string,
) *YoutubePlaylistPoller {
	poller := YoutubePlaylistPoller{
		youtubeApiKey:     youtubeApiKey,
		youtubePlaylistId: youtubePlaylistId,
		videoUrls:         videoUrls,
		prevCheckIdMap:    make(map[string]struct{}),
	}

	return &poller
}
