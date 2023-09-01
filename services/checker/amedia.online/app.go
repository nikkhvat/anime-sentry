package animegoorg_check

import (
	"anime-bot-schedule/models"
	"anime-bot-schedule/pkg/database"
	"anime-bot-schedule/pkg/message"
	parsing "anime-bot-schedule/services/parser/amedia.online"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Check(bot *tgbotapi.BotAPI, anime models.Anime) {
	db := database.GetDB()

	resp, err := parsing.Fetch(anime.URL)
	if err != nil {
		log.Printf("error fetching anime data: %s", err)
		return
	}

	if resp.AddedEpisode != anime.LastReleasedEpisode {
		var subscribers []models.Subscriber
		db.Where("anime_id = ?", anime.ID).Find(&subscribers)
		for _, subscriber := range subscribers {
			text := fmt.Sprintf("%s \n\nВышла новая серия на телеэкранах японии: %s", resp.Title, resp.AddedEpisode)

			msg := message.NewMessage{
				Text:        text,
				Photo:       resp.Poster,
				UserId:      subscriber.TelegramID,
				Link:        anime.URL,
				AnimeId:     anime.ID,
				LinkTitle:   "Перейти",
				DeletePrev:  true,
				Unsubscribe: true,
			}

			msg.Send(bot)
		}
		anime.LastReleasedEpisode = resp.AddedEpisode
		db.Save(&anime)
	}
}
