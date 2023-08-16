package main

import (
	env "anime-bot-schedule/pkg/env"
	"anime-bot-schedule/pkg/message"
	"anime-bot-schedule/repositories"
	"log"
	"regexp"
	"time"

	database "anime-bot-schedule/pkg/database"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"

	amediaonline "anime-bot-schedule/services/amedia.online"
	animegoorg "anime-bot-schedule/services/animego.org"
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

	updates, err := bot.GetUpdatesChan(updateConfig)

	// pattern := `^https://animego.org/anime/.*$`
	// regexp, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
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
		Text:   "Добро пожаловать в бот Anime Schedule!\n\nВам нужно прислать ссылку на аниме и я буду уведомлять вас о выходе новых аниме\n\nСайты которые поддерживаются на данный момент:\n - animego.org\n - amedia.online",
	}

	msg.Send(bot)
}

func handleUpdate(db *gorm.DB, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// * If service is animego.org

	animeGOregexp, _ := regexp.Compile(animegoorg.LINK_PATTERN)
	amediaOnline, _ := regexp.Compile(amediaonline.LINK_PATTERN)
	if animeGOregexp.MatchString(update.Message.Text) {
		msg := animegoorg.Handle(db, update)
		msg.UserId = update.Message.Chat.ID
		msg.Send(bot)

	} else if amediaOnline.MatchString(update.Message.Text) {
		msg := amediaonline.Handle(db, update)
		msg.UserId = update.Message.Chat.ID
		msg.Send(bot)

	} else {
		msg := message.NewMessage{
			UserId: update.Message.Chat.ID,
			Text:   "Не похоже что это ссылка на аниме.\nМы поддерживаем сервисы:\n\n- https://animego.org\n- https://amedia.online",
		}

		msg.Send(bot)
	}
}
