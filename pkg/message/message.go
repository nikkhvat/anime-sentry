package message

import (
	repositoriesmessage "anime-bot-schedule/repositories/message"
	"fmt"
	"log"

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
	DeletePrev  bool
}

func (msg NewMessage) Send(bot *tgbotapi.BotAPI) tgbotapi.Message {

	isLink := msg.Link != "" && msg.LinkTitle != ""

	isUnsubscribe := msg.Unsubscribe && msg.AnimeId != 0

	isRemovePrevMessage := msg.DeletePrev && msg.AnimeId != 0

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

		sentMessage, _ := bot.Send(newMsg)

		if isRemovePrevMessage {
			deletePrevMessage(bot, msg.UserId, msg.AnimeId, sentMessage.MessageID)
		}

		return sentMessage

	} else {
		newMsg := tgbotapi.NewMessage(msg.UserId, msg.Text)

		if !emptyKeyboard {
			newMsg.ReplyMarkup = keyboard
		}

		sentMessage, _ := bot.Send(newMsg)

		if isRemovePrevMessage {
			deletePrevMessage(bot, msg.UserId, msg.AnimeId, sentMessage.MessageID)
		}

		return sentMessage
	}
}

func deletePrevMessage(bot *tgbotapi.BotAPI, userId int64, animeId uint, mewMessageId int) error {

	prevMessageId, err := repositoriesmessage.GetLastMessage(animeId, userId)

	if prevMessageId != 0 {
		deletePrevMsg := tgbotapi.DeleteMessageConfig{
			ChatID:    userId,
			MessageID: int(prevMessageId),
		}

		_, err = bot.Request(deletePrevMsg)

		if err != nil {
			log.Printf("Failed to delete message: %s", err)
		}
	}

	if err != nil {
		log.Println(err)
	}

	repositoriesmessage.UpdateLastMessage(animeId, userId, mewMessageId)

	return nil
}
