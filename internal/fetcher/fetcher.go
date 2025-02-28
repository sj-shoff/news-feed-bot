package fetcher

import (
	"context"
	"news-feed-bot/internal/model"
	"time"
)

type ArticleStorage interface {
	Store(ctx context.Context, article model.Article) error
}

type SourcesProvider interface {
	Sources(ctx context.Context) ([]model.Source, error)
}

type Source interface {
	Id() int64
	Name() string
	Fetch(ctx context.Context) ([]model.Item, error)
}

type Fetcher struct {
	articles       ArticleStorage
	sources        SourcesProvider
	fetchInterval  time.Duration
	filterKeywords []string
}

func New(
	articleStorage ArticleStorage,
	sourcesProvider SourcesProvider,
	fetchInterval time.Duration,
	filterKeywords []string,
) *Fetcher {
	return &Fetcher{
		articles:       articleStorage,
		sources:        sourcesProvider,
		fetchInterval:  fetchInterval,
		filterKeywords: filterKeywords,
	}
}
