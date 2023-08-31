package fouranimeis

import (
	"anime-bot-schedule/models"
	"anime-bot-schedule/pkg/message"
	parsing "anime-bot-schedule/services/parser/4anime.is"
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

	if resp.LastEpisode != anime.LastReleasedEpisode {
		var subscribers []models.Subscriber
		db.Where("anime_id = ?", anime.ID).Find(&subscribers)
		for _, subscriber := range subscribers {
			text := fmt.Sprintf("%s \n\nA new series has been released %s", resp.Title, resp.LastEpisode)

			msg := message.NewMessage{
				Text:        text,
				Photo:       resp.Poster,
				UserId:      subscriber.TelegramID,
				Link:        resp.LastEpisodeLink,
				AnimeId:     anime.ID,
				LinkTitle:   "Open",
				Unsubscribe: true,
			}

			msg.Send(bot)
		}
		anime.LastReleasedEpisode = resp.LastEpisode
		db.Save(&anime)
	}
}
