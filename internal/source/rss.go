package source

import (
	"context"
	"github.com/SlyMarbo/rss"
	_ "github.com/samber/lo"
	"telegram-bot/internal/model"
)

type RSSSource struct {
	URL        string
	SourceID   int64
	SourceName int64
}

func NewRSSSourceFromModel(m model.Source) RSSSource {
	return RSSSource{
		URL:        m.FeedUrl,
		SourceID:   m.ID,
		SourceName: m.Name,
	}
}

func (s RSSSource) loadFeed(ctx context.Context, url string) (*rss.Feed, error) {
	var (
		feedCh = make(chan *rss.Feed)
		errCh  = make(chan error)
	)
	go func() {
		feed, err := rss.Fetch(url)
		if err != nil {
			errCh <- err
			return
		}

		feedCh <- feed
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case feed := <-feedCh:
		return feed, nil

	}
}
