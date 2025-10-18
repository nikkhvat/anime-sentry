package subscriber

import (
	"anime-sentry/models"
	"anime-sentry/parsers"
	"anime-sentry/pkg/localization"
	"anime-sentry/pkg/message"
	"anime-sentry/pkg/telegram"
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
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

	regexpAnimego, _ := regexp.Compile(`^https://animego.me/anime/.*$`)

	if regexpAnimego.MatchString(link) {
		return parsers.FetchAnimeGo(link)
	}

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

func (c *call) FollowAnime(ctx context.Context, command string, user models.User) error {

	var (
		messageAnimeSaved = localization.Localize(user.LanguageCode, "anime_saved")
	)

	tgbot := telegram.GetBot()

	parts := strings.Split(command, "_")

	if len(parts) != 2 {
		return errors.New("incorrect command")
	}

	animeId64, _ := strconv.ParseUint(parts[1], 10, 64)
	animeId := uint(animeId64)

	err := c.db.SubscribeToAnime(ctx, animeId, user.ID)

	if err != nil {
		if err.Error() == "you are already subscribed to this anime" {
			messageAlreadyTracking := localization.Localize(user.LanguageCode, "already_tracking")
			msg := message.NewMessage{Text: messageAlreadyTracking, UserId: user.ID}
			msg.Send(tgbot, user)

			return nil
		}

		messageUnknownError := localization.Localize(user.LanguageCode, "unknown_error")
		msg := message.NewMessage{Text: messageUnknownError, UserId: user.ID}
		msg.Send(tgbot, user)

		return nil
	}

	anime, err := c.db.GetAnimeById(ctx, animeId)

	if err != nil {
		log.Println(err)
	}

	data, err := getAnimeDataByUrl(anime.URL)

	if err != nil {
		messageUnknownError := localization.Localize(user.LanguageCode, "unknown_error")
		msg := message.NewMessage{Text: messageUnknownError, UserId: user.ID}
		msg.Send(tgbot, user)
	}

	_, nextEpisode := getLastAndNextEpisode(data.Episodes)

	messageWithEpisode := fmt.Sprintf(messageAnimeSaved, nextEpisode.Number, nextEpisode.Date)

	messageText := fmt.Sprintf("%s\n%s", escapeMarkdownV2(*data.Title), messageWithEpisode)

	newMsg := message.NewMessage{
		Text:         messageText,
		Unsubscribe:  true,
		AnimeId:      animeId,
		UserId:       user.ID,
		IsMarkdownV2: true,
	}

	newMsg.Send(tgbot, user)

	return nil
}
