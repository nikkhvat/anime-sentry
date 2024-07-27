package subscriber

import (
	"anime-sentry/models"
	"anime-sentry/repositories"
	"anime-sentry/services"
	"context"
)

type call struct {
	db repositories.DB
}

func (c *call) GetSubscriberByAnimeId(ctx context.Context, animeId uint) ([]models.User, error) {
	return c.db.GetSubscriberByAnimeId(ctx, animeId)
}

func (c *call) SubscribeToAnime(ctx context.Context, telegramID int64, url string, name string, image string, lastReleasedEpisode string, dubbings string) (*uint, error) {
	return c.db.SubscribeToAnime(ctx, telegramID, url, name, image, lastReleasedEpisode, dubbings)
}

func (c *call) UnsubscribeFromAnimeUpdates(ctx context.Context, animeId uint, userId int64) error {
	return c.db.UnsubscribeFromAnimeUpdates(ctx, animeId, userId)
}

func New(db repositories.DB) services.Subscriber {
	return &call{db: db}
}
