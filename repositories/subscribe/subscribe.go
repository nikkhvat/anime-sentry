package repositories_subscribe

import (
	"anime-bot-schedule/models"
	"anime-bot-schedule/pkg/database"
	"errors"

	"gorm.io/gorm"
)

func SubscribeToAnime(telegramID int64, url, name, image, lastReleasedEpisode string) (*uint, error) {
	db := database.GetDB()

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
			return nil, err
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
			return nil, err
		}
	} else {
		return nil, errors.New("you are already subscribed to this anime")
	}

	return &anime.ID, nil
}
