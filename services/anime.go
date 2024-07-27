package services

import (
	"anime-sentry/pkg/message"
	"context"
)

type Anime interface {
	SaveAnime(ctx context.Context, link string, userId int64) message.NewMessage
	// CheckAnime(ctx context.Context, link string, userId int64) message.NewMessage
	CheckAnimeStatus(ctx context.Context) error
}
