package bot

import (
	"context"
	"strings"

	"anime-sentry/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Callback(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	user := models.User{
		ID:        update.CallbackQuery.From.ID,
		UserName:  update.CallbackQuery.From.UserName,
		FirstName: update.CallbackQuery.From.FirstName,
		LastName:  update.CallbackQuery.From.LastName,
	}

	var err error

	language, err := h.user.Language(ctx, user.ID)

	if err != nil {
		user.LanguageCode = update.Message.From.LanguageCode
	} else {
		user.LanguageCode = *language
	}

	switch update.CallbackQuery.Data {
	case "en":
		user.LanguageCode = "en"
		h.user.ChooseLanguage(ctx, user)
	case "ru":
		user.LanguageCode = "ru"
		h.user.ChooseLanguage(ctx, user)
	}

	if strings.Contains(update.CallbackQuery.Data, "unsub") {
		h.subscriber.UnsubscribeFromAnimeUpdates(ctx, update.CallbackQuery.Data, update.CallbackQuery.Message, user)
	}

	if strings.Contains(update.CallbackQuery.Data, "follow") {
		h.subscriber.FollowAnime(ctx, update.CallbackQuery.Data, user)
	}
}
