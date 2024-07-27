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
	GetAnimeList(ctx context.Context) ([]models.Anime, error)
	UpdateLastEpisode(ctx context.Context, animeId uint, lastEpisode string) error
}

type Subscriber interface {
	SubscribeToAnime(ctx context.Context, telegramID int64, url, name, image, lastReleasedEpisode string, dubbings string) (*uint, error)
	GetSubscriberByAnimeId(ctx context.Context, animeId uint) ([]models.User, error)
	UnsubscribeFromAnimeUpdates(ctx context.Context, animeId uint, userId int64) error
}

type User interface {
	AddNewUser(ctx context.Context, user models.User) error
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	IsExist(ctx context.Context, user models.User) bool
	GetUserLanguage(ctx context.Context, user models.User) (*string, error)
	SetUserLanguage(ctx context.Context, user models.User) error
}

type Message interface {
	UpdateLastMessage(ctx context.Context, animeId uint, userId int64, lastMessageId int) error
	GetLastMessage(ctx context.Context, animeId uint, userId int64) (int64, error)
}
