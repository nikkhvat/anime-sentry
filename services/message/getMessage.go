package message

import "context"

func (c *call) GetLastMessage(ctx context.Context, animeId uint, userId int64) (int64, error) {
	return c.db.GetLastMessage(ctx, animeId, userId)
}
