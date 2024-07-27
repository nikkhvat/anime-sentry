package telegram

import (
	"anime-sentry/pkg/env"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var once sync.Once

func InitBot() *tgbotapi.BotAPI {
	log.Println("Bot Starting...")

	botToken := env.Get("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("authorized on account %s", bot.Self.UserName)

	return bot
}

func GetBot() *tgbotapi.BotAPI {
	once.Do(func() {
		bot = InitBot()
	})

	return bot
}

func GetUpdates() tgbotapi.UpdatesChannel {
	botInctance := GetBot()

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	return botInctance.GetUpdatesChan(updateConfig)
}
