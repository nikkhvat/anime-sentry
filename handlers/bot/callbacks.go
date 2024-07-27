package bot

import (
	"context"
	"log"
	"strconv"
	"strings"

	"anime-sentry/models"
	"anime-sentry/pkg/localization"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Callback(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	user := &models.User{
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
		h.user.ChooseLanguage(ctx, *user)
	case "ru":
		user.LanguageCode = "ru"
		h.user.ChooseLanguage(ctx, *user)
	}

	if strings.Contains(update.CallbackQuery.Data, "unsub") {
		handleUnsubscribe(h, ctx, tgbot, update, *user)
	}
}

// ! ПЕРЕНЕСТИ В СЕРВИС
func handleUnsubscribe(h *handler, ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update, user models.User) {
	callbackData := update.CallbackQuery.Data
	parts := strings.Split(callbackData, "_")

	if len(parts) != 3 {
		return
	}

	userId, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse UserId: %s", err)
		return
	}

	AnimeId, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		log.Fatalf("failed to parse AnimeId: %s", err)
		return
	}

	AnimeIdUint := uint(AnimeId)

	h.subscriber.UnsubscribeFromAnimeUpdates(ctx, AnimeIdUint, userId)

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := tgbot.Request(callback); err != nil {
		log.Println(err)
	}

	deleteMsg := tgbotapi.DeleteMessageConfig{
		ChatID:    update.CallbackQuery.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.MessageID,
	}
	_, err = tgbot.Request(deleteMsg)
	if err != nil {
		log.Printf("failed to delete message: %s", err)
	}

	messageTextUnsubscribe := localization.Localize(user.LanguageCode, "unsubscribe")

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messageTextUnsubscribe)
	if _, err := tgbot.Send(msg); err != nil {
		log.Println(err)
	}
}
