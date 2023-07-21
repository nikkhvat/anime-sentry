package repositories

import (
	"anime-bot-schedule/models"
	"errors"

	"gorm.io/gorm"
)

func AddSubscribeAnime(db *gorm.DB, telegramID int64, url, name, image, lastReleasedEpisode string) error {
	var anime models.Anime

	if err := db.Where("url = ?", url).First(&anime).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			anime = models.Anime{
				URL:                 url,
				Name:                name,
				Image:               image,
				LastReleasedEpisode: lastReleasedEpisode,
				IsSeasonOver:        false,
			}
			db.Create(&anime)
		} else {
			return err
		}
	}

	var subscriber models.Subscriber
	if err := db.Where("telegram_id = ? AND anime_id = ?", telegramID, anime.ID).First(&subscriber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			subscriber = models.Subscriber{
				TelegramID: telegramID,
				AnimeID:    anime.ID,
			}
			db.Create(&subscriber)
		} else {
			return err
		}
	} else {
		return errors.New("you are already subscribed to this anime")
	}

	return nil
}
