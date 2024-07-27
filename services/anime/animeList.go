package anime

import (
	"anime-sentry/models"
	"context"
)

func (c *call) GetAnimeList(ctx context.Context) ([]models.Anime, error) {
	return c.db.GetAnimeList(ctx)
}
