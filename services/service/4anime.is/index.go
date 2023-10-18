package fouranimeis

import (
	"anime-bot-schedule/pkg/message"
	repositories_subscribe "anime-bot-schedule/repositories/subscribe"
	"fmt"

	parsing "anime-bot-schedule/services/parser/4anime.is"
)

var LINK_PATTERN = `^https://4anime.is/.*$`
var LINK = `4anime.is`
var LANG = "en"

func Handle(userId int64, text string) message.NewMessage {
	data, err := parsing.Fetch(text)

	if err != nil {
		msg := message.NewMessage{
			Text: "An unknown error has occurred :(",
		}

		return msg
	}

	if len(data.Title) == 0 {
		msg := message.NewMessage{
			Text:  "We have not found such an anime",
			Photo: "https://animego.org/animego/images/404.gif",
		}

		return msg
	}

	animeId, err := repositories_subscribe.SubscribeToAnime(userId, text, data.Title, *&data.Poster, data.LastEpisode)

	if err != nil {
		if err.Error() == "you are already subscribed to this anime" {
			msg := message.NewMessage{
				Text: "Are you already tracking this anime!",
			}

			return msg
		}

		msg := message.NewMessage{
			Text: "An unknown error has occurred :(",
		}
		return msg
	}

	messageText := fmt.Sprintf("%s\n\nanime is saved, you will receive notifications about new episodes",
		data.Title)

	newMsg := message.NewMessage{
		Text:        messageText,
		Photo:       data.Poster,
		Unsubscribe: true,
		AnimeId:     *animeId,
	}

	return newMsg
}
