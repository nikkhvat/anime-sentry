package anime

import (
	"anime-sentry/pkg/message"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	parsers "anime-sentry/parsers"

	localization "anime-sentry/pkg/localization"
)

func escapeMarkdownV2(text string) string {
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}

func getAnimeDataByUrl(link string) (*parsers.AnimeResponse, error) {
	if len(link) == 0 {
		return nil, errors.New("link cannot by empty")
	}

	// regexp4anime, _ := regexp.Compile(`^https://4anime.gg/.*$`)
	// regexpAmedia, _ := regexp.Compile(`^https://amedia.site/.*$`)
	regexpAnimego, _ := regexp.Compile(`^https://animego.org/anime/.*$`)
	// regexpAnimevost, _ := regexp.Compile(`^https://animevost.org/tip/tv/.*$`)

	// * FOR 4ANIME.GG
	// if regexp4anime.MatchString(link) {
	// 	return parsers.FetchAnimeGo(link)
	// }

	// * FOR ANIMEDIA SITE
	// if regexpAmedia.MatchString(link) {
	// 	return parsers.FetchAnimeGo(link)
	// }

	// * FOR ANIMEGO.ORG
	if regexpAnimego.MatchString(link) {
		return parsers.FetchAnimeGo(link)
	}

	// * FOR ANIMEVOST
	// if regexpAnimevost.MatchString(link) {
	// 	return parsers.FetchAnimeGo(link)
	// }

	return nil, errors.New("incorrect link")
}

func getLastAndNextEpisode(episodes []parsers.Episode) (parsers.Episode, parsers.Episode) {
	var lastReleased parsers.Episode
	var nextEpisode parsers.Episode

	foundReleased := false

	for i := len(episodes) - 1; i >= 0; i-- {
		if episodes[i].Released && !foundReleased {
			lastReleased = episodes[i]
			foundReleased = true
		} else if foundReleased && !episodes[i].Released {
			nextEpisode = episodes[i]
			break
		}
	}

	return lastReleased, nextEpisode
}

func (c *call) CheckAnime(ctx context.Context, link string, userId int64) message.NewMessage {
	user, err := c.db.GetUserByID(ctx, userId)

	if err != nil {
		messageUnknownError := localization.Localize(user.LanguageCode, "unknown_error")
		msg := message.NewMessage{Text: messageUnknownError, UserId: userId}
		return msg
	}

	data, err := getAnimeDataByUrl(link)

	if err != nil {
		messageUnknownError := localization.Localize(user.LanguageCode, "unknown_error")
		msg := message.NewMessage{Text: messageUnknownError, UserId: userId}
		return msg
	}

	if len(*data.Title) == 0 {
		messageNotFound := localization.Localize(user.LanguageCode, "not_found")
		msg := message.NewMessage{
			Text:   messageNotFound,
			Photo:  "https://animego.org/animego/images/404.gif",
			UserId: userId,
		}

		return msg
	}

	lastReleased, nextEpisode := getLastAndNextEpisode(data.Episodes)

	joinedDubbings := strings.Join(data.Dubbings, ", ")

	animeId, err := c.db.SubscribeToAnime(
		ctx,
		userId,
		link,
		*data.Title,
		*data.Image,
		lastReleased.Number,
		joinedDubbings,
	)

	if err != nil {
		if err.Error() == "you are already subscribed to this anime" {
			messageAlreadyTracking := localization.Localize(user.LanguageCode, "already_tracking")
			msg := message.NewMessage{Text: messageAlreadyTracking, UserId: userId}

			return msg
		}

		messageUnknownError := localization.Localize(user.LanguageCode, "unknown_error")
		msg := message.NewMessage{Text: messageUnknownError, UserId: userId}

		return msg
	}

	messageAnimeSaved := localization.Localize(user.LanguageCode, "anime_saved")

	messageWithEpisode := fmt.Sprintf(messageAnimeSaved, nextEpisode.Number, nextEpisode.Date)

	messageText := fmt.Sprintf("%s\n%s", escapeMarkdownV2(*data.Title), messageWithEpisode)

	newMsg := message.NewMessage{
		Text:         messageText,
		Photo:        *data.Image,
		Unsubscribe:  true,
		AnimeId:      *animeId,
		UserId:       userId,
		IsMarkdownV2: true,
	}

	return newMsg
}
