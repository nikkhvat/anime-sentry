package services

import (
	"context"

	"anime-sentry/models"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Subscriber interface {
	FollowAnime(ctx context.Context, command string, user models.User) error
	UnsubscribeFromAnimeUpdates(ctx context.Context, command string, message *tgBotApi.Message, user models.User) error
}
