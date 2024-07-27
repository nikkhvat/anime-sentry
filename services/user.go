package services

import (
	"context"

	"anime-sentry/models"
)

type User interface {
	AddNewUser(ctx context.Context, user models.User) error
	ChooseLanguage(ctx context.Context, user models.User) error
	IsExist(ctx context.Context, id int64) bool
	Language(ctx context.Context, id int64) (*string, error)
}
