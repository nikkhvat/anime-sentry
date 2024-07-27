package services

import (
	"anime-sentry/models"
	"anime-sentry/pkg/message"
	"context"
)

type Anime interface {
	GetAnimeList(ctx context.Context) ([]models.Anime, error)
	UpdateLastEpisode(ctx context.Context, animeId uint, lastEpisode string) error

	CheckAnime(ctx context.Context, link string, userId int64) message.NewMessage
	CheckAnimeStatus(ctx context.Context) error
}
