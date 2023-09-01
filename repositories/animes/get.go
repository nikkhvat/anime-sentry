package repositories_animes

import (
	"anime-bot-schedule/models"
	"anime-bot-schedule/pkg/database"
	"log"
)

func GetAnimes() ([]models.Anime, error) {
	db := database.GetDB()

	var animes []models.Anime
	if err := db.Find(&animes).Error; err != nil {
		log.Printf("error retrieving anime from DB: %s", err)
		return nil, err
	}

	return animes, nil
}
