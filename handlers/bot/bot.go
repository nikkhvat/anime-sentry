package bot

import (
	"anime-sentry/handlers"
	"anime-sentry/pkg/telegram"
	"anime-sentry/repositories"
	"context"

	anime "anime-sentry/services/anime"
	message "anime-sentry/services/message"
	subscriber "anime-sentry/services/subscriber"
	user "anime-sentry/services/user"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type bot struct {
	handler *handler
}

func New(db repositories.DB) handlers.Bot {
	anime := anime.New(db)
	message := message.New(db)
	subscriber := subscriber.New(db)
	user := user.New(db)

	handler := handler{
		anime:      anime,
		message:    message,
		subscriber: subscriber,
		user:       user,
	}

	return &bot{handler: &handler}
}

func (b *bot) Start(ctx context.Context) error {
	bot := telegram.GetBot()

	u := tgBotApi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message != nil && update.Message.IsCommand() {
			go b.handler.Command(ctx, bot, update)
		} else if update.CallbackQuery != nil {
			go b.handler.Callback(ctx, bot, update)
		} else if update.Message != nil {
			go b.handler.Message(ctx, bot, update)
		} else {
			continue
		}
	}

	return nil
}
