package message

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NewMessage struct {
	Text        string
	Photo       string
	UserId      int64
	AnimeId     uint
	Link        string
	LinkTitle   string
	Unsubscribe bool
}

func (msg NewMessage) Send(bot *tgbotapi.BotAPI) {

	isLink := msg.Link != "" && msg.LinkTitle != ""
	isUnsubscribe := msg.Unsubscribe && msg.AnimeId != 0

	emptyKeyboard := !isLink && !msg.Unsubscribe

	var keyboard tgbotapi.InlineKeyboardMarkup

	if isLink && isUnsubscribe {
		unsubButtonData := fmt.Sprintf("unsub_%d_%d", msg.UserId, msg.AnimeId)

		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(msg.LinkTitle, msg.Link),
				tgbotapi.NewInlineKeyboardButtonData("Отписаться", unsubButtonData),
			),
		)
	} else if isUnsubscribe {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Отписаться", "unsubscribe"),
			),
		)
	} else if isLink {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(msg.LinkTitle, msg.Link),
			),
		)
	}

	if msg.Photo != "" {
		file := tgbotapi.FileURL(msg.Photo)
		newMsg := tgbotapi.NewPhoto(msg.UserId, file)
		newMsg.Caption = msg.Text

		if !emptyKeyboard {
			newMsg.ReplyMarkup = keyboard
		}

		_, _ = bot.Send(newMsg)
	} else {
		newMsg := tgbotapi.NewMessage(msg.UserId, msg.Text)

		if !emptyKeyboard {
			newMsg.ReplyMarkup = keyboard
		}

		_, _ = bot.Send(newMsg)
	}
}
