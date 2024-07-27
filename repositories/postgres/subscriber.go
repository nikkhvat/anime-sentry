package postgres

import (
	models "anime-sentry/models"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

func (p *postgres) GetSubscriberByAnimeId(ctx context.Context, animeId uint) ([]models.User, error) {
	if p.db == nil {
		return nil, errors.New("no database connection")
	}

	var subscribers []models.User

	result := p.db.Table("subscribers").Select("subscribers.*, users.id, users.first_name, users.last_name, users.user_name, COALESCE(users.language_code, 'en') as language_code").
		Joins("LEFT JOIN users ON users.id = subscribers.telegram_id").
		Where("subscribers.anime_id = ?", animeId).
		Scan(&subscribers)

	if result.Error != nil {
		return nil, result.Error
	}

	return subscribers, nil
}

func (p *postgres) SubscribeToAnime(ctx context.Context, animeId uint, userId int64) error {
	var subscriber models.Subscriber
	if err := p.db.Where("telegram_id = ? AND anime_id = ?", userId, animeId).First(&subscriber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			subscriber = models.Subscriber{
				TelegramID: userId,
				AnimeID:    animeId,
			}
			p.db.Create(&subscriber)
		} else {
			return err
		}
	} else {
		return errors.New("you are already subscribed to this anime")
	}

	return nil
}

func (p *postgres) UnsubscribeFromAnimeUpdates(ctx context.Context, animeId uint, userId int64) error {
	var subscription models.Subscriber

	result := p.db.Where("anime_id = ? AND telegram_id = ?", animeId, userId).First(&subscription)

	if result.Error != nil {
		log.Printf("error finding subscription: %s", result.Error)
		return result.Error
	}
	p.db.Unscoped().Delete(&subscription)

	var subscribersCount int64
	p.db.Model(&models.Subscriber{}).Where("anime_id = ?", animeId).Count(&subscribersCount)

	return nil
}
