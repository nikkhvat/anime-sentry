package repositories

import (
	"anime-bot-schedule/models"
	"anime-bot-schedule/parsing"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
)

func CheckAnimeStatus(db *gorm.DB, bot *tgbotapi.BotAPI) {
	var animes []models.Anime
	if err := db.Find(&animes).Error; err != nil {
		log.Printf("Error retrieving anime from DB: %s", err)
		return
	}

	for _, anime := range animes {
		resp, err := parsing.AnimeGOFetch(anime.URL)
		if err != nil {
			log.Printf("Error fetching anime data: %s", err)
			continue
		}

		var lastEpisod parsing.Episod

		if resp.Episods[0].Relized {
			lastEpisod = resp.Episods[0]
		} else if resp.Episods[1].Relized {
			lastEpisod = resp.Episods[1]
		} else if resp.Episods[2].Relized {
			lastEpisod = resp.Episods[2]
		}

		if lastEpisod.Number != anime.LastReleasedEpisode {
			var subscribers []models.Subscriber
			db.Where("anime_id = ?", anime.ID).Find(&subscribers)
			for _, subscriber := range subscribers {
				text := fmt.Sprintf("%s \n\nВышла новая серия на телеэкранах японии %s (%s)\n%s", *resp.Title, lastEpisod.Number, lastEpisod.Title, anime.URL)

				if resp.Image != nil {
					msg := tgbotapi.NewPhotoShare(subscriber.TelegramID, *resp.Image)
					msg.Caption = text
					_, _ = bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(subscriber.TelegramID, text)
					_, _ = bot.Send(msg)
				}
			}
			anime.LastReleasedEpisode = lastEpisod.Number
			db.Save(&anime)
		}
	}
}
