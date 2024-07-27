package services

import (
	"context"

	"anime-sentry/models"
)

type Subscriber interface {
	SubscribeToAnime(ctx context.Context, telegramID int64, url, name, image, lastReleasedEpisode string, dubbings string) (*uint, error)
	GetSubscriberByAnimeId(ctx context.Context, animeId uint) ([]models.User, error)
	UnsubscribeFromAnimeUpdates(ctx context.Context, animeId uint, userId int64) error
}
