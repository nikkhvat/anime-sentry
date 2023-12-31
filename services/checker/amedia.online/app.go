package animegoorg_check

import (
	"anime-bot-schedule/models"
	"anime-bot-schedule/pkg/localization"
	"anime-bot-schedule/pkg/message"
	repositories_animes "anime-bot-schedule/repositories/animes"
	repositories_subscribe "anime-bot-schedule/repositories/subscribe"
	parsing "anime-bot-schedule/services/parser/amedia.online"
	"fmt"
	"log"
)

func Check(anime models.Anime) {
	resp, err := parsing.Fetch(anime.URL)
	if err != nil {
		log.Printf("error fetching anime data: %s", err)
		return
	}

	if resp.AddedEpisode != anime.LastReleasedEpisode {

		subscribers, err := repositories_subscribe.GetByAnime(anime.ID)

		if err != nil {
			return
		}

		for _, subscriber := range subscribers {
			messageNewSeries := localization.Localize(subscriber.LanguageCode, "new_series")
			text := fmt.Sprintf("%s \n%s", resp.Title, messageNewSeries)

			openButtonText := localization.Localize(subscriber.LanguageCode, "open_link")

			msg := message.NewMessage{
				Text:        text,
				Photo:       resp.Poster,
				UserId:      subscriber.TelegramID,
				Link:        anime.URL,
				AnimeId:     anime.ID,
				LinkTitle:   openButtonText,
				DeletePrev:  true,
				Unsubscribe: true,
			}

			msg.Send()
		}

		repositories_animes.UpdateLastEpisod(anime.ID, resp.AddedEpisode)
	}
}
