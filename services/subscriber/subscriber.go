package subscriber

import (
	"anime-sentry/repositories"
	"anime-sentry/services"
)

type call struct {
	db repositories.DB
}

func New(db repositories.DB) services.Subscriber {
	return &call{db: db}
}
