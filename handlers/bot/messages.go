package bot

import (
	"context"
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

	regexp4anime, _ := regexp.Compile(`^https://4anime.gg/.*$`)
	regexpAmedia, _ := regexp.Compile(`^https://amedia.site/.*$`)
	regexpAnimego, _ := regexp.Compile(`^https://animego.org/anime/.*$`)
	regexpAnimevost, _ := regexp.Compile(`^https://animevost.org/tip/tv/.*$`)

	if regexpAnimego.MatchString(link) || regexp4anime.MatchString(link) ||
		regexpAmedia.MatchString(link) || regexpAnimevost.MatchString(link) {
		msg := h.anime.CheckAnime(ctx, update.Message.Text, user.ID)
		msg.Send(tgbot, *user)
		return
	}

	result := generateAnimeSitesMessage(localization.Localize(user.LanguageCode, "invalid_link"))

	msg := message.NewMessage{
		UserId: user.ID,
		Text:   result,
	}

	msg.Send(tgbot, *user)
}
