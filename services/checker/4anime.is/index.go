package fouranimeis_check

import (
	"anime-bot-schedule/models"
	"anime-bot-schedule/pkg/message"
	repositories_animes "anime-bot-schedule/repositories/animes"
	repositories_subscribe "anime-bot-schedule/repositories/subscribe"
	parsing "anime-bot-schedule/services/parser/4anime.is"
	"fmt"
	"log"
)

func Check(anime models.Anime) {
	resp, err := parsing.Fetch(anime.URL)
	if err != nil {
		log.Printf("error fetching anime data: %s", err)
		return
	}

	if resp.LastEpisode != anime.LastReleasedEpisode {
		subscribers, err := repositories_subscribe.GetByAnime(anime.ID)

		if err != nil {
			return
		}

		for _, subscriber := range subscribers {
			text := fmt.Sprintf("%s \n\nA new series has been released %s", resp.Title, resp.LastEpisode)

			msg := message.NewMessage{
				Text:        text,
				Photo:       resp.Poster,
				UserId:      subscriber.TelegramID,
				Link:        resp.LastEpisodeLink,
				AnimeId:     anime.ID,
				LinkTitle:   "Open",
				DeletePrev:  true,
				Unsubscribe: true,
			}

			msg.Send()
		}
		repositories_animes.UpdateLastEpisod(anime.ID, resp.LastEpisode)
	}
}
