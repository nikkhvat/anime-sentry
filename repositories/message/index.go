package repositories_message

import (
	"anime-bot-schedule/models"
	"anime-bot-schedule/pkg/database"
	"errors"

	"gorm.io/gorm"
)

func GetLastMessage(animeId uint, userId int64) (int64, error) {
	db := database.GetDB()

	var subscriber models.Subscriber

	if err := db.Where("anime_id = ? AND telegram_id = ?", animeId, userId).First(&subscriber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("Subscriber not found")
		}
		return 0, err
	}

	return subscriber.LastMessage, nil
}

func UpdateLastMessage(animeId uint, userId int64, lastMessageId int) error {
	db := database.GetDB()

	var subscriber models.Subscriber

	if err := db.Where("anime_id = ? AND telegram_id = ?", animeId, userId).First(&subscriber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Subscriber not found")
		}
		return err
	}

	subscriber.LastMessage = int64(lastMessageId)
	if err := db.Save(&subscriber).Error; err != nil {
		return err
	}

	return nil
}
