package fouranimeis

import (
	"anime-bot-schedule/models"
	parsing "anime-bot-schedule/services/parser/4anime.is"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
)

func Check(db *gorm.DB, bot *tgbotapi.BotAPI, anime models.Anime) {
	resp, err := parsing.Fetch(anime.URL)
	if err != nil {
		log.Printf("error fetching anime data: %s", err)
		return
	}

	if resp.LastEpisode != anime.LastReleasedEpisode {
		var subscribers []models.Subscriber
		db.Where("anime_id = ?", anime.ID).Find(&subscribers)
		for _, subscriber := range subscribers {
			text := fmt.Sprintf("%s \n\nA new series has been released %s \n%s", resp.Title, resp.LastEpisode, resp.LastEpisodeLink)

			if resp.Poster != "" {
				msg := tgbotapi.NewPhotoShare(subscriber.TelegramID, resp.Poster)
				msg.Caption = text
				_, _ = bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(subscriber.TelegramID, text)
				_, _ = bot.Send(msg)
			}
		}
		anime.LastReleasedEpisode = resp.LastEpisode
		db.Save(&anime)
	}
}
