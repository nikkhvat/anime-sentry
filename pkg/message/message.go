package message

import (
	"anime-sentry/models"
	"anime-sentry/pkg/localization"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NewMessage struct {
	Text         string
	Photo        string
	UserId       int64
	AnimeId      uint
	Link         string
	LinkTitle    string
	Unsubscribe  bool
	DeletePrev   bool
	IsMarkdownV2 bool
}

// * Send message in telegram
func (msg NewMessage) Send(tgbot *tgbotapi.BotAPI, user models.User) *tgbotapi.Message {
	if msg.UserId == 0 {
		return nil
	}

	isLink := msg.Link != "" && msg.LinkTitle != ""

	isUnsubscribe := msg.Unsubscribe && msg.AnimeId != 0

	emptyKeyboard := !isLink && !msg.Unsubscribe

	var keyboard tgbotapi.InlineKeyboardMarkup

	if isLink && isUnsubscribe {
		unsubButtonData := fmt.Sprintf("unsub_%d_%d", msg.UserId, msg.AnimeId)

		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(msg.LinkTitle, msg.Link),

				tgbotapi.NewInlineKeyboardButtonData(
					localization.Localize(user.LanguageCode, "unsubscribe_button"),
					unsubButtonData),
			),
		)
	} else if isUnsubscribe {
		unsubButtonData := fmt.Sprintf("unsub_%d_%d", msg.UserId, msg.AnimeId)

		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(

				tgbotapi.NewInlineKeyboardButtonData(
					localization.Localize(user.LanguageCode, "unsubscribe_button"),
					unsubButtonData),
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

		if msg.IsMarkdownV2 {
			newMsg.ParseMode = "MarkdownV2"
		}
		sentMessage, err := tgbot.Send(newMsg)

		if err != nil {
			log.Println(err)
		}

		return &sentMessage

	} else {
		newMsg := tgbotapi.NewMessage(msg.UserId, msg.Text)

		if !emptyKeyboard {
			newMsg.ReplyMarkup = keyboard
		}

		if msg.IsMarkdownV2 {
			newMsg.ParseMode = "MarkdownV2"
		}
		sentMessage, err := tgbot.Send(newMsg)

		if err != nil {
			log.Println(err)
		}

		return &sentMessage
	}
}
