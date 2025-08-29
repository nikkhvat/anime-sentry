package postgres

import (
	models "anime-sentry/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

func (p *postgres) AddNewUser(ctx context.Context, user models.User) error {
	result := p.db.FirstOrCreate(&user, models.User{ID: user.ID})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (p *postgres) GetUserAnimeList(ctx context.Context, user models.User) ([]models.Anime, error) {
	var animes []models.Anime

	err := p.db.WithContext(ctx).
		Preload("Anime").
		Where("telegram_id = ?", user.ID).
		Find(&[]models.Subscriber{}).Error
	if err != nil {
		return nil, err
	}

	var subscribers []models.Subscriber
	result := p.db.WithContext(ctx).
		Preload("Anime").
		Where("telegram_id = ?", user.ID).
		Find(&subscribers)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, subscriber := range subscribers {
		animes = append(animes, subscriber.Anime)
	}

	return animes, nil
}

func (p *postgres) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User

	result := p.db.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (p *postgres) IsExist(ctx context.Context, user models.User) bool {

	var userFromDB models.User

	result := p.db.First(&userFromDB, user.ID)

	return result.Error == nil
}

func (p *postgres) GetUserLanguage(ctx context.Context, user models.User) (*string, error) {
	var userFromDB models.User

	result := p.db.First(&userFromDB, user.ID)

	if result.Error != nil {
		return nil, result.Error
	}

	return &userFromDB.LanguageCode, nil
}

func (p *postgres) SetUserLanguage(ctx context.Context, user models.User) error {
	result := p.db.Model(&models.User{}).Where("id = ?", user.ID).Update("language_code", user.LanguageCode)

	return result.Error
}
