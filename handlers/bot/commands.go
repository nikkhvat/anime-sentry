package bot

import (
	"anime-sentry/pkg/localization"
	"anime-sentry/pkg/message"
	"context"
	"fmt"
	"strings"

	"anime-sentry/models"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Command(ctx context.Context, tgbot *tgBotApi.BotAPI, update tgBotApi.Update) {
	user := models.User{
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

	switch update.Message.Command() {
	case "start":
		if !h.user.IsExist(ctx, user.ID) {
			h.user.AddNewUser(ctx, user)
		}

		messageText := localization.Localize(user.LanguageCode, "welcome")

		result := generateAnimeSitesMessage(messageText)

		msg := message.NewMessage{
			UserId: user.ID,
			Text:   result,
		}

		msg.Send(tgbot, user)
	}
}

func generateAnimeSitesMessage(message string) string {
	var siteLinks []string

	for _, site := range []string{"4anime.gg", "animego.org", "amedia.site", "animevost.org"} {
		siteLinks = append(siteLinks, site)
	}

	formattedSites := "- " + strings.Join(siteLinks, "\n- ")

	fullMessage := fmt.Sprintf("%s\n\n%s", message, formattedSites)
	return fullMessage
}
