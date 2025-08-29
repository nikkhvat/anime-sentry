package user

import (
	"anime-sentry/models"
	"anime-sentry/repositories"
	"anime-sentry/services"
	"context"
)

type call struct {
	db repositories.DB
}

func (c *call) AddNewUser(ctx context.Context, user models.User) error {
	return c.db.AddNewUser(ctx, user)
}

func (c *call) ChooseLanguage(ctx context.Context, user models.User) error {
	return c.db.SetUserLanguage(ctx, user)
}

func (c *call) GetUserAnimeList(ctx context.Context, user models.User) ([]models.Anime, error) {
	return c.db.GetUserAnimeList(ctx, user)
}

func (c *call) IsExist(ctx context.Context, id int64) bool {
	user := models.User{
		ID: id,
	}

	isExist := c.db.IsExist(ctx, user)

	return isExist
}

func (c *call) Language(ctx context.Context, id int64) (*string, error) {
	defaultLanguage := "en"

	user := models.User{
		ID: id,
	}

	lang, err := c.db.GetUserLanguage(ctx, user)

	if err != nil {
		return &defaultLanguage, err
	}

	return lang, nil
}

func New(db repositories.DB) services.User {
	return &call{db: db}
}
