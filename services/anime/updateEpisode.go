package anime

import "context"

func (c *call) UpdateLastEpisode(ctx context.Context, animeId uint, lastEpisode string) error {
	return c.db.UpdateLastEpisode(ctx, animeId, lastEpisode)
}
