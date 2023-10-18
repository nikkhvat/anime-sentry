package users_repository

import (
	"anime-bot-schedule/models"
	"anime-bot-schedule/pkg/database"
	"errors"

	"gorm.io/gorm"
)

func AddUser(user models.User) error {
	db := database.GetDB()

	result := db.FirstOrCreate(&user, models.User{ID: user.ID})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetUserByID(id int64) (*models.User, error) {
	db := database.GetDB()

	var user models.User

	result := db.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
