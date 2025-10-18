package bot

import (
	"context"
	"log"
	"regexp"

	"anime-sentry/models"
	"anime-sentry/pkg/localization"
	"anime-sentry/pkg/message"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AnimeService struct {
	Link        string
	LinkPattern string
}

func (h *handler) Message(ctx context.Context, tgbot *tgBotApi.BotAPI, update tgBotApi.Update) {
	user := &models.User{
		ID:        update.Message.From.ID,
		UserName:  update.Message.From.UserName,
		FirstName: update.Message.From.FirstName,
		LastName:  update.Message.From.LastName,
	}

	var err error
	language, err := h.user.Language(ctx, user.ID)

	if err != nil {
		user.LanguageCode = update.Message.From.LanguageCode
	} else {
		user.LanguageCode = *language
	}

	link := update.Message.Text

	regexpAnimego, _ := regexp.Compile(`^https://animego.me/anime/.*$`)

	if regexpAnimego.MatchString(link) {
		msg := h.anime.SaveAnime(ctx, update.Message.Text, user.ID)
		msg.Send(tgbot, *user)
		return
	}

	if update.Message.Text == localization.Localize(user.LanguageCode, "change_language") {
		languageMsg := tgBotApi.NewMessage(user.ID, localization.Localize(user.LanguageCode, "choose_language"))

		button1 := tgBotApi.NewInlineKeyboardButtonData("Русский 🇷🇺", "ru")
		button2 := tgBotApi.NewInlineKeyboardButtonData("English 🇺🇸", "en")

		row := tgBotApi.NewInlineKeyboardRow(button1, button2)
		languageKeyboard := tgBotApi.NewInlineKeyboardMarkup(row)

		languageMsg.ReplyMarkup = languageKeyboard
		_, err := tgbot.Send(languageMsg)

		if err != nil {
			log.Println(err)
		}

		return
	}

	result := generateAnimeSitesMessage(localization.Localize(user.LanguageCode, "invalid_link"))

	msg := message.NewMessage{
		UserId: user.ID,
		Text:   result,
	}

	msg.Send(tgbot, *user)
}
