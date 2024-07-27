package anime

import (
	localization "anime-sentry/pkg/localization"
	"anime-sentry/pkg/message"
	"anime-sentry/pkg/telegram"
	"context"
	"errors"
	"fmt"
)

func (c *call) CheckAnimeStatus(ctx context.Context) error {
	bot := telegram.GetBot()

	list, err := c.db.GetAnimeList(ctx)

	if err != nil {
		return errors.New("cannot get list of anime")
	}

	for _, anime := range list {
		data, err := getAnimeDataByUrl(anime.URL)

		if err != nil {
			// TODO: Записывать в логе что что то пошло не так
			continue
		}

		lastReleased, _ := getLastAndNextEpisode(data.Episodes)

		if lastReleased.Number != anime.LastReleasedEpisode && len(data.Dubbings) > 0 {
			users, err := c.db.GetSubscriberByAnimeId(ctx, anime.ID)

			if err != nil {
				continue
			}

			for _, user := range users {
				messageNewSeries := localization.Localize(user.LanguageCode, "new_series")
				openLinkButtonText := localization.Localize(user.LanguageCode, "open_link")
				text := fmt.Sprintf("%s \n%s", anime.Name, messageNewSeries)

				msg := message.NewMessage{
					Text:        text,
					Photo:       anime.Image,
					UserId:      user.ID,
					Link:        anime.URL,
					AnimeId:     anime.ID,
					LinkTitle:   openLinkButtonText,
					DeletePrev:  true,
					Unsubscribe: true,
				}

				msg.Send(bot, user)
			}

			c.db.UpdateLastEpisode(ctx, anime.ID, lastReleased.Number)
		}
	}

	return nil
}
