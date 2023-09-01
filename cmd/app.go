package main

import (
	message "anime-bot-schedule/pkg/message"
	telegram "anime-bot-schedule/pkg/telegram"

	repositories_check "anime-bot-schedule/repositories/check"
	repositories_subscribe "anime-bot-schedule/repositories/subscribe"

	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	gocron "github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	fouranimeis "anime-bot-schedule/services/service/4anime.is"
	amediaonline "anime-bot-schedule/services/service/amedia.online"
	animegoorg "anime-bot-schedule/services/service/animego.org"
	animevostorg "anime-bot-schedule/services/service/animevost.org"
)

func main() {
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
		if update.CallbackQuery != nil {
			handleUnsub(update)
			continue
		}

		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/start" {
			go startBot(update.Message.Chat.ID)
			continue
		}

		go handleUpdate(update.Message.Chat.ID, update.Message.Text)
	}
}

func handleUnsub(update tgbotapi.Update) {
	bot := telegram.GetBot()

	callbackData := update.CallbackQuery.Data
	parts := strings.Split(callbackData, "_")

	if len(parts) != 3 {
		log.Fatalf("Unexpected callback data: %s", callbackData)
		return
	}

	userId, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse UserId: %s", err)
		return
	}

	AnimeId, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse AnimeId: %s", err)
		return
	}

	AnimeIdUint := uint(AnimeId)

	repositories_subscribe.Unsubscribe(AnimeIdUint, userId)

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := bot.Request(callback); err != nil {
		panic(err)
	}

	deleteMsg := tgbotapi.DeleteMessageConfig{
		ChatID:    update.CallbackQuery.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.MessageID,
	}
	_, err = bot.Request(deleteMsg)
	if err != nil {
		log.Printf("failed to delete message: %s", err)
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы отписались от этого аниме!")
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func startBot(userId int64) {
	msg := message.NewMessage{
		UserId: userId,
		Text:   "Добро пожаловать в бот Anime Schedule!\n\nВам нужно прислать ссылку на аниме и я буду уведомлять вас о выходе новых аниме\n\nСайты которые поддерживаются на данный момент:\n- animego.org\n- amedia.online\n- animevost.org\n- 4anime.is",
	}

	msg.Send()
}

func handleUpdate(userId int64, messageText string) {
	// * If service is animego.org

	animeGOregexp, _ := regexp.Compile(animegoorg.LINK_PATTERN)
	amediaOnline, _ := regexp.Compile(amediaonline.LINK_PATTERN)
	animevostOrg, _ := regexp.Compile(animevostorg.LINK_PATTERN)
	fouranimeIs, _ := regexp.Compile(fouranimeis.LINK_PATTERN)

	if animeGOregexp.MatchString(messageText) {
		msg := animegoorg.Handle(userId, messageText)
		msg.UserId = userId
		msg.Send()
	} else if amediaOnline.MatchString(messageText) {
		msg := amediaonline.Handle(userId, messageText)
		msg.UserId = userId
		msg.Send()
	} else if animevostOrg.MatchString(messageText) {
		msg := animevostorg.Handle(userId, messageText)
		msg.UserId = userId
		msg.Send()
	} else if fouranimeIs.MatchString(messageText) {
		msg := fouranimeis.Handle(userId, messageText)
		msg.UserId = userId
		msg.Send()
	} else {
		msg := message.NewMessage{
			UserId: userId,
			Text:   "Не похоже что это ссылка на аниме.\nМы поддерживаем сервисы:\n\n- animego.org\n- amedia.online\n- animevost.org\n- 4anime.is",
		}

		msg.Send()
	}
}
