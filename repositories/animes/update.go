package repositories_animes

import (
	"anime-bot-schedule/models"
	"anime-bot-schedule/pkg/database"
)

func UpdateLastEpisod(animeId uint, lastEpisode string) error {
	db := database.GetDB()

	var anime models.Anime

	if err := db.First(&anime, animeId).Error; err != nil {
		return err
	}

	anime.LastReleasedEpisode = lastEpisode
	if err := db.Save(&anime).Error; err != nil {
		return err
	}

	return nil
}
