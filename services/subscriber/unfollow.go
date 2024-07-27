package subscriber

import (
	"anime-sentry/models"
	"anime-sentry/pkg/localization"
	"anime-sentry/pkg/telegram"
	"context"
	"errors"
	"strconv"
	"strings"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *call) UnsubscribeFromAnimeUpdates(ctx context.Context, command string, message *tgBotApi.Message, user models.User) error {
	tgbot := telegram.GetBot()

	parts := strings.Split(command, "_")

	if len(parts) != 3 {
		return errors.New("incorrect command")
	}

	userId, _ := strconv.ParseInt(parts[1], 10, 64)
	animeId64, _ := strconv.ParseUint(parts[2], 10, 64)
	animeId := uint(animeId64)

	c.db.UnsubscribeFromAnimeUpdates(ctx, animeId, userId)

	deleteMsg := tgBotApi.DeleteMessageConfig{
		ChatID:    message.Chat.ID,
		MessageID: message.MessageID,
	}

	_, err := tgbot.Request(deleteMsg)
	if err != nil {
		return err
	}

	messageTextUnsubscribe := localization.Localize(user.LanguageCode, "unsubscribe")

	msg := tgBotApi.NewMessage(message.Chat.ID, messageTextUnsubscribe)
	if _, err := tgbot.Send(msg); err != nil {
		return err
	}

	return nil
}
