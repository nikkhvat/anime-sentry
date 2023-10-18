package repositories_subscribe

import (
	"anime-bot-schedule/models"
	"anime-bot-schedule/pkg/database"
	"errors"
)

type SubscriberInfo struct {
	models.Subscriber
	FirstName    string
	LastName     string
	UserName     string
	LanguageCode string `gorm:"default:'en'"`
}

func GetByAnime(animeId uint) ([]SubscriberInfo, error) {
	db := database.GetDB()

	if db == nil {
		return nil, errors.New("no database connection")
	}

	var subscribersInfo []SubscriberInfo

	result := db.Table("subscribers").Select("subscribers.*, users.first_name, users.last_name, users.user_name, COALESCE(users.language_code, 'en') as language_code").
		Joins("LEFT JOIN users ON users.id = subscribers.telegram_id").
		Where("subscribers.anime_id = ?", animeId).
		Scan(&subscribersInfo)

	if result.Error != nil {
		return nil, result.Error
	}

	return subscribersInfo, nil
}
