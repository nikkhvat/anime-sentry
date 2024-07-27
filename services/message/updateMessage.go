package message

import "context"

func (c *call) UpdateLastMessage(ctx context.Context, animeId uint, userId int64, lastMessageId int) error {
	return c.db.UpdateLastMessage(ctx, animeId, userId, lastMessageId)
}
