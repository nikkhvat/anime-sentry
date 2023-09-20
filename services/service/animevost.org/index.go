package animevostorg

import (
	"anime-bot-schedule/pkg/message"
	repositories_subscribe "anime-bot-schedule/repositories/subscribe"
	"fmt"

	parsing "anime-bot-schedule/services/parser/animevost.org"
)

var LINK_PATTERN = `^https://animevost.org/tip/tv/.*$`

func Handle(userId int64, text string) message.NewMessage {
	data, err := parsing.Fetch(text)

	if err != nil {
		msg := message.NewMessage{
			Text: "Произошла ошибка :(",
		}

		return msg
	}

	if len(data.Title) == 0 {
		msg := message.NewMessage{
			Text:  "Мы не нашли такого аниме",
			Photo: "https://animego.org/animego/images/404.gif",
		}

		return msg
	}

	animeId, err := repositories_subscribe.SubscribeToAnime(userId, text, data.Title, *&data.Poster, data.AddedEpisode)

	if err != nil {
		if err.Error() == "you are already subscribed to this anime" {
			msg := message.NewMessage{
				Text: "Вы уже подписанны на это аниме!",
			}

			return msg
		}

		msg := message.NewMessage{
			Text: "Произошла неизвестная ошибка :(",
		}
		return msg
	}

	messageText := fmt.Sprintf("%s\n\nАниме сохраненно, вы будете получать уведомления когда выйдут новые серии. \n\n%s выйдет %s.",
		data.Title, data.NextEpisode, data.NextEpisodeDate)

	newMsg := message.NewMessage{
		Text:        messageText,
		Unsubscribe: true,
		AnimeId:     *animeId,
	}

	if data.Poster != "" {
		newMsg.Photo = data.Poster
	}

	return newMsg
}
