package animevostorg

import (
	"anime-bot-schedule/pkg/message"
	repositories_subscribe "anime-bot-schedule/repositories/subscribe"
	"fmt"

	localization "anime-bot-schedule/pkg/localization"
	parsing "anime-bot-schedule/services/parser/animevost.org"
)

var LINK_PATTERN = `^https://animevost.org/tip/tv/.*$`
var LINK = `animevost.org`
var LANG = "ru"

func Handle(userId int64, text string, lang string) message.NewMessage {
	data, err := parsing.Fetch(text)

	if err != nil {
		messageUnknownError := localization.Localize(lang, "unknown_error")
		msg := message.NewMessage{Text: messageUnknownError}

		return msg
	}

	if len(data.Title) == 0 {
		messageNotFound := localization.Localize(lang, "not_found")
		msg := message.NewMessage{
			Text:  messageNotFound,
			Photo: "https://animego.org/animego/images/404.gif",
		}

		return msg
	}

	animeId, err := repositories_subscribe.SubscribeToAnime(userId, text, data.Title, *&data.Poster, data.AddedEpisode)

	if err != nil {
		if err.Error() == "you are already subscribed to this anime" {
			messageAlreadyTracking := localization.Localize(lang, "already_tracking")
			msg := message.NewMessage{
				Text: messageAlreadyTracking,
			}

			return msg
		}

		messageUnknownError := localization.Localize(lang, "unknown_error")
		msg := message.NewMessage{Text: messageUnknownError}
		return msg
	}

	messageAnimeSaved := localization.Localize(lang, "anime_saved")
	messageText := fmt.Sprintf("%s\n%s", data.Title, messageAnimeSaved)

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
