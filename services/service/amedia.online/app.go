package amediaonline

import (
	"anime-bot-schedule/pkg/message"
	"anime-bot-schedule/repositories"
	parsing "anime-bot-schedule/services/parser/amedia.online"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

var LINK_PATTERN = `^https://amedia.online/.*$`

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
		data.Title, *&data.Poster, data.AddedEpisode)

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
		Text: messageText,
	}

	if data.Poster != "" {
		newMsg.Photo = data.Poster
	}

	return newMsg
}
