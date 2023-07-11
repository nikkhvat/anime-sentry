package main

import (
	parsing "anime-bot-schedule/parsing"
	env "anime-bot-schedule/pkg/env"
	"anime-bot-schedule/repositories"
	"fmt"
	"log"
	"regexp"

	database "anime-bot-schedule/pkg/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
)

func main() {
	db := database.InitDB()

	botToken := env.Get("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)

	pattern := `^https://animego.org/anime/.*$`
	regexp, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go handleUpdate(db, bot, update, regexp)
	}
}

func handleUpdate(db *gorm.DB, bot *tgbotapi.BotAPI, update tgbotapi.Update, regexp *regexp.Regexp) {
	if regexp.MatchString(update.Message.Text) {

		resp, err := parsing.AnimeGOFetch(update.Message.Text)

		if err != nil {
			log.Panic(err)
		}

		if len(*resp.Title) == 0 {
			msg := tgbotapi.NewPhotoShare(update.Message.Chat.ID, "https://animego.org/animego/images/404.gif")
			msg.Caption = "Мы не нашли такого аниме"
			_, _ = bot.Send(msg)
			return
		}

		var lastEpisod parsing.Episod

		if !resp.Episods[0].Relized && resp.Episods[1].Relized {
			lastEpisod = resp.Episods[0]
		} else if !resp.Episods[1].Relized && resp.Episods[2].Relized {
			lastEpisod = resp.Episods[1]
		} else {
			lastEpisod = resp.Episods[2]
		}

		err = repositories.AddSubscribeAnime(db, update.Message.Chat.ID, update.Message.Text,
			*resp.Title, *resp.Image, lastEpisod.Number)

		if err != nil {
			log.Println(err)
		}

		message := fmt.Sprintf("%s\n\nАниме сохраненно, вы будете получать уведомления когда выйдут новые серии. \n\n%s (%s) выйдет %s.",
			*resp.Title, lastEpisod.Number, lastEpisod.Title, lastEpisod.Date)

		if resp.Image != nil && *resp.Image != "" {
			msg := tgbotapi.NewPhotoShare(update.Message.Chat.ID, *resp.Image)
			msg.Caption = message
			_, _ = bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			_, _ = bot.Send(msg)
		}
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не похоже что это ссылка на animego.org :(")
		_, _ = bot.Send(msg)
	}
}
