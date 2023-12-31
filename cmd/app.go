package main

import (
	"anime-bot-schedule/models"
	message "anime-bot-schedule/pkg/message"
	telegram "anime-bot-schedule/pkg/telegram"
	"fmt"

	repositories_check "anime-bot-schedule/repositories/check"
	repositories_subscribe "anime-bot-schedule/repositories/subscribe"

	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	gocron "github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"

	fouranimeis "anime-bot-schedule/services/service/4anime.is"
	amediaonline "anime-bot-schedule/services/service/amedia.online"
	animegoorg "anime-bot-schedule/services/service/animego.org"
	animevostorg "anime-bot-schedule/services/service/animevost.org"

	users_repository "anime-bot-schedule/repositories/users"

	localization "anime-bot-schedule/pkg/localization"
)

type AnimeService struct {
	Link        string
	LinkPattern string
	Lang        string
	Handle      func(userId int64, text string, lang string) message.NewMessage
}

var SERVICES = []AnimeService{
	{
		Link:        fouranimeis.LINK,
		LinkPattern: fouranimeis.LINK_PATTERN,
		Lang:        fouranimeis.LANG,
		Handle:      fouranimeis.Handle,
	},
	{
		Link:        amediaonline.LINK,
		LinkPattern: amediaonline.LINK_PATTERN,
		Lang:        amediaonline.LANG,
		Handle:      amediaonline.Handle,
	},
	{
		Link:        animegoorg.LINK,
		LinkPattern: animegoorg.LINK_PATTERN,
		Lang:        animegoorg.LANG,
		Handle:      animegoorg.Handle,
	},
	{
		Link:        animevostorg.LINK,
		LinkPattern: animevostorg.LINK_PATTERN,
		Lang:        animevostorg.LANG,
		Handle:      animevostorg.Handle,
	},
}

func main() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	// * Launch a goroutine for regular status checks (every 30 minutes)
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(30).Minute().Do(func() {
		repositories_check.CheckAnimeStatus()
	})

	if err != nil {
		log.Fatal(err)
	}

	s.StartAsync()

	updates := telegram.GetUpdates()
	for update := range updates {

		if update.Message != nil {
			users_repository.AddUser(models.User{
				ID:           update.Message.From.ID,
				FirstName:    update.Message.From.FirstName,
				LastName:     update.Message.From.LastName,
				UserName:     update.Message.From.UserName,
				LanguageCode: update.Message.From.LanguageCode,
			})
		}

		if update.CallbackQuery != nil {
			handleUnsub(update)
			continue
		}

		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/start" {
			go startBot(update.Message.Chat.ID, update.Message.From.LanguageCode)
			continue
		}

		go handleUpdate(update.Message.Chat.ID, update.Message.Text, update.Message.From.LanguageCode)
	}
}

func GenerateAnimeSitesMessage(message string, sites []AnimeService) string {
	var siteLinks []string

	for _, site := range sites {
		siteLinks = append(siteLinks, site.Link)
	}

	formattedSites := "- " + strings.Join(siteLinks, "\n- ")

	fullMessage := fmt.Sprintf("%s\n\n%s", message, formattedSites)
	return fullMessage
}

func handleUnsub(update tgbotapi.Update) {
	bot := telegram.GetBot()

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

	repositories_subscribe.Unsubscribe(AnimeIdUint, userId)

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := bot.Request(callback); err != nil {
		log.Println(err)
	}

	deleteMsg := tgbotapi.DeleteMessageConfig{
		ChatID:    update.CallbackQuery.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.MessageID,
	}
	_, err = bot.Request(deleteMsg)
	if err != nil {
		log.Printf("failed to delete message: %s", err)
	}

	messageTextUnsubscribe := localization.Localize(update.CallbackQuery.From.LanguageCode, "unsubscribe")

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messageTextUnsubscribe)
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func startBot(userId int64, lang string) {
	messageText := localization.Localize(lang, "welcome")

	result := GenerateAnimeSitesMessage(messageText, SERVICES)

	msg := message.NewMessage{
		UserId: userId,
		Text:   result,
	}

	msg.Send()
}

func handleUpdate(userId int64, messageText string, lang string) {
	matchCount := 0

	for _, service := range SERVICES {
		regexp, _ := regexp.Compile(service.LinkPattern)

		if regexp.MatchString(messageText) {
			msg := service.Handle(userId, messageText, lang)
			msg.UserId = userId
			matchCount++
			msg.Send()
		}
	}

	if matchCount == 0 {
		messageText := localization.Localize(lang, "invalid_link")

		result := GenerateAnimeSitesMessage(messageText, SERVICES)

		msg := message.NewMessage{
			UserId: userId,
			Text:   result,
		}

		msg.Send()
	}
}
