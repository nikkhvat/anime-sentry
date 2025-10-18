package bot

import (
	"anime-sentry/pkg/localization"
	"context"
	"fmt"
	"log"
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

	var (
		changeLanguageButtonText = localization.Localize(user.LanguageCode, "change_language")
	)

	var generalKeyboard = tgBotApi.NewReplyKeyboard(
		tgBotApi.NewKeyboardButtonRow(
			tgBotApi.NewKeyboardButton(changeLanguageButtonText),
		),
	)

	switch update.Message.Command() {
	case "start":
		if !h.user.IsExist(ctx, user.ID) {
			h.user.AddNewUser(ctx, user)
		}

		messageText := localization.Localize(user.LanguageCode, "welcome")

		result := generateAnimeSitesMessage(messageText)

		msg := tgBotApi.NewMessage(user.ID, result)
		msg.ReplyMarkup = generalKeyboard

		tgbot.Send(msg)
	case "language":
		onChangeLanguage(&user, tgbot)
	case "list":
		anime, err := h.user.GetUserAnimeList(ctx, user)

		if err != nil {
			log.Println(err)
			return
		}

		onListAnime(&user, tgbot, anime)
	}
}

func generateAnimeSitesMessage(message string) string {
	var siteLinks []string

	for _, site := range []string{"animego.me"} {
		siteLinks = append(siteLinks, site)
	}

	formattedSites := "- " + strings.Join(siteLinks, "\n- ")

	fullMessage := fmt.Sprintf("%s\n\n%s", message, formattedSites)
	return fullMessage
}

func onListAnime(user *models.User, tgbot *tgBotApi.BotAPI, anime []models.Anime) {
	var (
		langMessage = localization.Localize(user.LanguageCode, "anime_list")
	)

	var animeList []string
	for _, a := range anime {
		animeList = append(animeList, fmt.Sprintf("- %s", a.Name))
	}

	messageText := fmt.Sprintf("%s\n%s", langMessage, strings.Join(animeList, "\n"))

	languageMsg := tgBotApi.NewMessage(user.ID, messageText)

	_, err := tgbot.Send(languageMsg)

	if err != nil {
		log.Println(err)
	}
}
