package repositories_subscribe

import (
	"anime-bot-schedule/models"
	"anime-bot-schedule/pkg/database"
)

func GetByAnime(animeId uint) ([]models.Subscriber, error) {
	db := database.GetDB()

	var subscribers []models.Subscriber
	result := db.Where("anime_id = ?", animeId).Find(&subscribers)

	if result.Error != nil {
		return nil, result.Error
	}

	return subscribers, nil
}
