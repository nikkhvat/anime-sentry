package main

import (
	env "anime-bot-schedule/pkg/env"
	"anime-bot-schedule/pkg/message"
	"anime-bot-schedule/repositories"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	database "anime-bot-schedule/pkg/database"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"

	fouranimeis "anime-bot-schedule/services/service/4anime.is"
	amediaonline "anime-bot-schedule/services/service/amedia.online"
	animegoorg "anime-bot-schedule/services/service/animego.org"
	animevostorg "anime-bot-schedule/services/service/animevost.org"
)

func main() {

	db := database.InitDB()

	botToken := env.Get("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	// * Launch a goroutine for regular status checks (every 30 minutes)
	s := gocron.NewScheduler(time.UTC)

	_, err = s.Every(30).Minute().Do(func() {
		repositories.CheckAnimeStatus(db, bot)
	})

	if err != nil {
		log.Fatalf("Could not schedule job: %v", err)
	}

	s.StartAsync()

	// * Start a TG Bot
	log.Printf("authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {

		if update.CallbackQuery != nil {
			callbackData := update.CallbackQuery.Data
			parts := strings.Split(callbackData, "_")

			if len(parts) != 3 {
				log.Fatalf("Unexpected callback data: %s", callbackData)
				return
			}

			UserId, err := strconv.ParseInt(parts[1], 10, 64)
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
			log.Printf("UserId: %d, AnimeId: %d", UserId, AnimeIdUint)

			repositories.Unsubscribe(db, AnimeIdUint, UserId)

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
				log.Printf("Failed to delete message: %s", err)
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы отписались от этого аниме!")
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}

		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/start" {
			go startBot(bot, update)
			continue
		}

		go handleUpdate(db, bot, update)
	}
}

func startBot(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := message.NewMessage{
		UserId: update.Message.Chat.ID,
		Text:   "Добро пожаловать в бот Anime Schedule!\n\nВам нужно прислать ссылку на аниме и я буду уведомлять вас о выходе новых аниме\n\nСайты которые поддерживаются на данный момент:\n- animego.org\n- amedia.online\n- animevost.org\n- 4anime.is",
	}

	msg.Send(bot)
}

func handleUpdate(db *gorm.DB, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// * If service is animego.org

	animeGOregexp, _ := regexp.Compile(animegoorg.LINK_PATTERN)
	amediaOnline, _ := regexp.Compile(amediaonline.LINK_PATTERN)
	animevostOrg, _ := regexp.Compile(animevostorg.LINK_PATTERN)
	fouranimeIs, _ := regexp.Compile(fouranimeis.LINK_PATTERN)

	if animeGOregexp.MatchString(update.Message.Text) {
		msg := animegoorg.Handle(db, update)
		msg.UserId = update.Message.Chat.ID
		msg.Send(bot)
	} else if amediaOnline.MatchString(update.Message.Text) {
		msg := amediaonline.Handle(db, update)
		msg.UserId = update.Message.Chat.ID
		msg.Send(bot)
	} else if animevostOrg.MatchString(update.Message.Text) {
		msg := animevostorg.Handle(db, update)
		msg.UserId = update.Message.Chat.ID
		msg.Send(bot)
	} else if fouranimeIs.MatchString(update.Message.Text) {
		msg := fouranimeis.Handle(db, update)
		msg.UserId = update.Message.Chat.ID
		msg.Send(bot)
	} else {
		msg := message.NewMessage{
			UserId: update.Message.Chat.ID,
			Text:   "Не похоже что это ссылка на аниме.\nМы поддерживаем сервисы:\n\n- animego.org\n- amedia.online\n- animevost.org\n- 4anime.is",
		}

		msg.Send(bot)
	}
}
