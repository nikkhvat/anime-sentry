package anime

import (
	"anime-sentry/models"
	parsers "anime-sentry/parsers"
	localization "anime-sentry/pkg/localization"
	"anime-sentry/pkg/message"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

func getAnimeDataByUrl(link string) (*parsers.AnimeResponse, error) {
	if len(link) == 0 {
		return nil, errors.New("link cannot by empty")
	}

	regexpAnimego, _ := regexp.Compile(`^https://animego.me/anime/.*$`)

	if regexpAnimego.MatchString(link) {
		return parsers.FetchAnimeGo(link)
	}

	return nil, errors.New("incorrect link")
}

func (c *call) SaveAnime(ctx context.Context, link string, userId int64) message.NewMessage {
	user, err := c.db.GetUserByID(ctx, userId)

	var (
		messageUnknownError = localization.Localize(user.LanguageCode, "unknown_error")
		messageNotFound     = localization.Localize(user.LanguageCode, "not_found")
		followButtonText    = localization.Localize(user.LanguageCode, "follow")
	)

	if err != nil {
		msg := message.NewMessage{Text: messageUnknownError, UserId: userId}
		return msg
	}

	data, err := getAnimeDataByUrl(link)

	if err != nil {
		msg := message.NewMessage{Text: messageUnknownError, UserId: userId}
		return msg
	}

	if len(*data.Title) == 0 {
		msg := message.NewMessage{
			Text:   messageNotFound,
			Photo:  "https://animego.me/animego/images/404.gif",
			UserId: userId,
		}

		return msg
	}

	lastReleased, _ := getLastAndNextEpisode(data.Episodes)

	joinedDubbings := strings.Join(data.Dubbings, ", ")

	anime := models.Anime{
		URL:                 link,
		Name:                *data.Title,
		Image:               *data.Image,
		LastReleasedEpisode: lastReleased.Number,
		// TODO: Проверять вышел ли сезон или нет
		IsSeasonOver: false,
		Dubbings:     joinedDubbings,
	}

	id, err := c.db.SaveAnime(ctx, anime)

	// TODO: Смотреть какая ошибка, если аниме уже есть, просто отправить сообщение пользователю с кнопками
	if err != nil {
		return message.NewMessage{
			Text:   messageUnknownError,
			UserId: userId,
		}
	}

	unsubscribeButtonValue := fmt.Sprintf("follow_%d", id)

	return message.NewMessage{
		UserId:           userId,
		Text:             *data.Title,
		Photo:            *data.Image,
		IsCustomKeyboard: true,
		CustomKeyboard: tgBotApi.NewInlineKeyboardMarkup(
			tgBotApi.NewInlineKeyboardRow(
				tgBotApi.NewInlineKeyboardButtonData(followButtonText, unsubscribeButtonValue),
			),
		),
	}
}
