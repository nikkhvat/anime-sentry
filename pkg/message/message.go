package message

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type NewMessage struct {
	Text   string `json:"text"`
	Photo  string `json:"photo"`
	UserId int64  `json:"user_id"`
}

func (msg NewMessage) Send(bot *tgbotapi.BotAPI) {
	if msg.Photo != "" {
		newMsg := tgbotapi.NewPhotoShare(msg.UserId, msg.Photo)
		newMsg.Caption = msg.Text
		_, _ = bot.Send(newMsg)
	} else {
		msg := tgbotapi.NewMessage(msg.UserId, msg.Text)
		_, _ = bot.Send(msg)
	}
}
