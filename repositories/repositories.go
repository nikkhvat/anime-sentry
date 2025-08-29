package repositories

import (
	"anime-sentry/models"
	"context"
)

type DB interface {
	User
	Anime
	Subscriber
	Message
}

type Anime interface {
	SaveAnime(ctx context.Context, anime models.Anime) (uint, error)
	GetAnimeList(ctx context.Context) ([]models.Anime, error)
	GetAnimeById(ctx context.Context, animeId uint) (models.Anime, error)
	UpdateLastEpisode(ctx context.Context, animeId uint, lastEpisode string) error
}

type Subscriber interface {
	SubscribeToAnime(ctx context.Context, animeId uint, userId int64) error
	GetSubscriberByAnimeId(ctx context.Context, animeId uint) ([]models.User, error)
	UnsubscribeFromAnimeUpdates(ctx context.Context, animeId uint, userId int64) error
}

type User interface {
	AddNewUser(ctx context.Context, user models.User) error
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	IsExist(ctx context.Context, user models.User) bool
	GetUserLanguage(ctx context.Context, user models.User) (*string, error)
	SetUserLanguage(ctx context.Context, user models.User) error
	GetUserAnimeList(ctx context.Context, user models.User) ([]models.Anime, error)
}

type Message interface {
	UpdateLastMessage(ctx context.Context, animeId uint, userId int64, lastMessageId int) error
	GetLastMessage(ctx context.Context, animeId uint, userId int64) (int64, error)
}
