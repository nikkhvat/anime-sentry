package postgres

import (
	"context"
	"errors"
	"log"

	models "anime-sentry/models"

	"gorm.io/gorm"
)

func (p *postgres) GetAnimeList(ctx context.Context) ([]models.Anime, error) {
	var animeList []models.Anime

	if err := p.db.Find(&animeList).Error; err != nil {
		log.Printf("error retrieving anime from DB: %s", err)
		return nil, err
	}

	return animeList, nil
}

func (p *postgres) UpdateLastEpisode(ctx context.Context, animeId uint, lastEpisode string) error {
	var anime models.Anime

	if err := p.db.First(&anime, animeId).Error; err != nil {
		return err
	}

	anime.LastReleasedEpisode = lastEpisode
	if err := p.db.Save(&anime).Error; err != nil {
		return err
	}

	return nil
}

func (p *postgres) SaveAnime(ctx context.Context, anime models.Anime) (uint, error) {
	var existingAnime models.Anime

	result := p.db.Where("url = ?", anime.URL).First(&existingAnime)
	if result.Error == nil {
		return existingAnime.ID, nil
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, result.Error
	}

	createResult := p.db.Create(&anime)
	return anime.ID, createResult.Error
}

func (p *postgres) GetAnimeById(ctx context.Context, animeId uint) (models.Anime, error) {
	var anime models.Anime

	result := p.db.First(&anime, animeId)
	if result.Error != nil {
		return models.Anime{}, result.Error
	}

	return anime, nil
}
