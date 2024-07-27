package bot

import (
	"anime-sentry/services"
)

type handler struct {
	anime      services.Anime
	message    services.Message
	subscriber services.Subscriber
	user       services.User
}
