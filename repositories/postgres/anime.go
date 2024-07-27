package postgres

import (
	"context"
	"log"

	models "anime-sentry/models"
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
