package animegoorgcheck

import (
	"anime-bot-schedule/models"
	"anime-bot-schedule/pkg/message"
	parsing "anime-bot-schedule/services/parser/animego.org"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func Check(db *gorm.DB, bot *tgbotapi.BotAPI, anime models.Anime) {
	resp, err := parsing.Fetch(anime.URL)
	if err != nil {
		log.Printf("error fetching anime data: %s", err)
		return
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
			text := fmt.Sprintf("%s \n\nВышла новая серия на телеэкранах японии: %s (%s)", *resp.Title, lastEpisod.Number, lastEpisod.Title)

			msg := message.NewMessage{
				Text:        text,
				Photo:       *resp.Image,
				UserId:      subscriber.TelegramID,
				Link:        anime.URL,
				AnimeId:     anime.ID,
				LinkTitle:   "Перейти",
				Unsubscribe: true,
			}

			msg.Send(bot)
		}
		anime.LastReleasedEpisode = lastEpisod.Number
		db.Save(&anime)
	}
}
