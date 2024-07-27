package services

import (
	"context"
)

type Message interface {
	UpdateLastMessage(ctx context.Context, animeId uint, userId int64, lastMessageId int) error
	GetLastMessage(ctx context.Context, animeId uint, userId int64) (int64, error)
}
