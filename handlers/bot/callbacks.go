package bot

import (
	"context"
	"strings"

	"anime-sentry/models"
	"anime-sentry/pkg/localization"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Callback(ctx context.Context, tgbot *tgBotApi.BotAPI, update tgBotApi.Update) {
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
		changeLang(ctx, "en", &user, tgbot, h)
		return
	case "ru":
		changeLang(ctx, "ru", &user, tgbot, h)
		return
	}

	if strings.Contains(update.CallbackQuery.Data, "unsub") {
		h.subscriber.UnsubscribeFromAnimeUpdates(ctx, update.CallbackQuery.Data, update.CallbackQuery.Message, user)
	}

	if strings.Contains(update.CallbackQuery.Data, "follow") {
		h.subscriber.FollowAnime(ctx, update.CallbackQuery.Data, user)
	}
}

func changeLang(ctx context.Context, lang string, user *models.User, tgbot *tgBotApi.BotAPI, h *handler) {
	user.LanguageCode = lang

	err := h.user.ChooseLanguage(ctx, *user)

	if err != nil {
		msg := tgBotApi.NewMessage(user.ID, localization.Localize(user.LanguageCode, "unknown_error"))
		tgbot.Send(msg)
		return
	}

	msg := tgBotApi.NewMessage(user.ID, localization.Localize(user.LanguageCode, "language_successfully_changed"))
	var generalKeyboard = tgBotApi.NewReplyKeyboard(
		tgBotApi.NewKeyboardButtonRow(
			tgBotApi.NewKeyboardButton(localization.Localize(user.LanguageCode, "change_language")),
		),
	)
	msg.ReplyMarkup = generalKeyboard
	tgbot.Send(msg)
}
