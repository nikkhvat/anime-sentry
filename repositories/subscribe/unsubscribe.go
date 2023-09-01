package repositories_subscribe

import (
	"log"

	"anime-bot-schedule/models"
	"anime-bot-schedule/pkg/database"
)

func Unsubscribe(animeId uint, userId int64) error {
	db := database.GetDB()

	var subscription models.Subscriber
	result := db.Where("anime_id = ? AND telegram_id = ?", animeId, userId).First(&subscription)
	if result.Error != nil {
		log.Printf("Error finding subscription: %s", result.Error)
		return result.Error
	}
	db.Unscoped().Delete(&subscription)

	var subscribersCount int64
	db.Model(&models.Subscriber{}).Where("anime_id = ?", animeId).Count(&subscribersCount)

	if subscribersCount == 0 {
		var anime models.Anime
		result = db.First(&anime, animeId)
		if result.Error != nil {
			log.Printf("Error finding anime: %s", result.Error)
			return result.Error
		}
		db.Unscoped().Delete(&anime)
	}

	log.Printf("Successfully unsubscribed user %d from anime %d", userId, animeId)
	return nil
}
