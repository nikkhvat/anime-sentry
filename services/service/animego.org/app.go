package animegoorg

import (
	localization "anime-bot-schedule/pkg/localization"
	"anime-bot-schedule/pkg/message"
	repositories_subscribe "anime-bot-schedule/repositories/subscribe"
	parsing "anime-bot-schedule/services/parser/animego.org"
	"fmt"
)

var LINK_PATTERN = `^https://animego.org/anime/.*$`
var LINK = `animego.org`
var LANG = "ru"

func Handle(userId int64, text string, lang string) message.NewMessage {
	data, err := parsing.Fetch(text)

	if err != nil {
		messageUnknownError := localization.Localize(lang, "unknown_error")
		msg := message.NewMessage{Text: messageUnknownError}

		return msg
	}

	if len(*data.Title) == 0 {
		messageNotFound := localization.Localize(lang, "not_found")
		msg := message.NewMessage{
			Text:  messageNotFound,
			Photo: "https://animego.org/animego/images/404.gif",
		}

		return msg
	}

	var lastEpisod parsing.Episod

	if data.Episods[0].Relized {
		lastEpisod = data.Episods[0]
	} else if data.Episods[1].Relized {
		lastEpisod = data.Episods[1]
	} else if data.Episods[2].Relized {
		lastEpisod = data.Episods[2]
	}

	animeId, err := repositories_subscribe.SubscribeToAnime(userId, text, *data.Title, *data.Image, lastEpisod.Number)

	if err != nil {
		if err.Error() == "you are already subscribed to this anime" {
			messageAlreadyTracking := localization.Localize(lang, "already_tracking")
			msg := message.NewMessage{Text: messageAlreadyTracking}

			return msg
		}

		messageUnknownError := localization.Localize(lang, "unknown_error")
		msg := message.NewMessage{Text: messageUnknownError}
		return msg
	}

	messageAnimeSaved := localization.Localize(lang, "anime_saved")
	messageText := fmt.Sprintf("%s\n%s", lastEpisod.Title, messageAnimeSaved)

	newMsg := message.NewMessage{
		Text:        messageText,
		Unsubscribe: true,
		AnimeId:     *animeId,
	}

	if data.Image != nil && *data.Image != "" {
		newMsg.Photo = *data.Image
	}

	return newMsg
}
