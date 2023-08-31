package fouranimeis

import (
	"anime-bot-schedule/pkg/message"
	"anime-bot-schedule/repositories"
	"fmt"

	parsing "anime-bot-schedule/services/parser/4anime.is"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

var LINK_PATTERN = `^https://4anime.is/.*$`

func Handle(db *gorm.DB, update tgbotapi.Update) message.NewMessage {
	data, err := parsing.Fetch(update.Message.Text)

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

	err = repositories.SubscribeToAnime(db, update.Message.Chat.ID, update.Message.Text,
		data.Title, *&data.Poster, data.LastEpisode)

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

	messageText := fmt.Sprintf("%s\n\nanime is saved, you will receive notifications about new episodes",
		data.Title)

	newMsg := message.NewMessage{
		Text:  messageText,
		Photo: data.Poster,
	}

	return newMsg
}
