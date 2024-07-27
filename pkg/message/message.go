package message

import (
	"anime-sentry/models"
	"anime-sentry/pkg/localization"
	"fmt"
	"log"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NewMessage struct {
	Text             string
	Photo            string
	UserId           int64
	AnimeId          uint
	Link             string
	LinkTitle        string
	Unsubscribe      bool
	DeletePrev       bool
	IsMarkdownV2     bool
	IsCustomKeyboard bool
	CustomKeyboard   tgBotApi.InlineKeyboardMarkup
}

// * Send message in telegram
func (msg NewMessage) Send(tgbot *tgBotApi.BotAPI, user models.User) *tgBotApi.Message {

	var (
		unsubscribeButtonText  = localization.Localize(user.LanguageCode, "unsubscribe_button")
		unsubscribeButtonValue = fmt.Sprintf("unsub_%d_%d", msg.UserId, msg.AnimeId)
	)

	if msg.UserId == 0 {
		return nil
	}

	isLink := msg.Link != "" && msg.LinkTitle != ""

	isUnsubscribe := msg.Unsubscribe && msg.AnimeId != 0

	var keyboard tgBotApi.InlineKeyboardMarkup

	var buttons []tgBotApi.InlineKeyboardButton

	if isUnsubscribe {
		buttons = append(buttons, tgBotApi.NewInlineKeyboardButtonData(unsubscribeButtonText, unsubscribeButtonValue))
	}

	if isLink {
		buttons = append(buttons, tgBotApi.NewInlineKeyboardButtonURL(msg.LinkTitle, msg.Link))
	}

	if msg.IsCustomKeyboard {
		keyboard = msg.CustomKeyboard
	} else {
		keyboard = tgBotApi.NewInlineKeyboardMarkup(buttons)
	}

	isKeyboard := len(buttons) > 0 || msg.IsCustomKeyboard

	if msg.Photo != "" {
		return sendPhotoMessage(tgbot, msg, keyboard, isKeyboard)
	}

	return sendTextMessage(tgbot, msg, keyboard, isKeyboard)
}

func sendPhotoMessage(tgbot *tgBotApi.BotAPI, msg NewMessage, keyboard tgBotApi.InlineKeyboardMarkup, isKeyboard bool) *tgBotApi.Message {
	file := tgBotApi.FileURL(msg.Photo)
	newMsg := tgBotApi.NewPhoto(msg.UserId, file)
	newMsg.Caption = msg.Text

	if isKeyboard {
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

func sendTextMessage(tgbot *tgBotApi.BotAPI, msg NewMessage, keyboard tgBotApi.InlineKeyboardMarkup, isKeyboard bool) *tgBotApi.Message {
	newMsg := tgBotApi.NewMessage(msg.UserId, msg.Text)

	if isKeyboard {
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
