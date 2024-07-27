package postgres

import (
	models "anime-sentry/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

func (p *postgres) UpdateLastMessage(ctx context.Context, animeId uint, userId int64, lastMessageId int) error {
	var subscriber models.Subscriber

	if err := p.db.Where("anime_id = ? AND telegram_id = ?", animeId, userId).First(&subscriber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("subscriber not found")
		}
		return err
	}

	subscriber.LastMessage = int64(lastMessageId)
	if err := p.db.Save(&subscriber).Error; err != nil {
		return err
	}

	return nil
}

func (p *postgres) GetLastMessage(ctx context.Context, animeId uint, userId int64) (int64, error) {
	var subscriber models.Subscriber

	if err := p.db.Where("anime_id = ? AND telegram_id = ?", animeId, userId).First(&subscriber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("subscriber not found")
		}
		return 0, err
	}

	return subscriber.LastMessage, nil
}
