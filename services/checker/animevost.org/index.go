package animevostorgcheck

import (
	"anime-bot-schedule/models"
	parsing "anime-bot-schedule/services/parser/animevost.org"
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

	if resp.AddedEpisode != anime.LastReleasedEpisode {
		var subscribers []models.Subscriber
		db.Where("anime_id = ?", anime.ID).Find(&subscribers)
		for _, subscriber := range subscribers {
			text := fmt.Sprintf("%s \n\nВышла новая серия на телеэкранах японии %s \n%s", resp.Title, resp.AddedEpisode, anime.URL)

			if resp.Poster != "" {
				msg := tgbotapi.NewPhotoShare(subscriber.TelegramID, resp.Poster)
				msg.Caption = text
				_, _ = bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(subscriber.TelegramID, text)
				_, _ = bot.Send(msg)
			}
		}
		anime.LastReleasedEpisode = resp.AddedEpisode
		db.Save(&anime)
	}
}
