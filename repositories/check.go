package repositories

import (
	amediaonlinecheck "anime-bot-schedule/checker/amedia.online"
	animegoorgcheck "anime-bot-schedule/checker/animego.org"
	"anime-bot-schedule/models"
	"log"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
)

func CheckAnimeStatus(db *gorm.DB, bot *tgbotapi.BotAPI) {
	var animes []models.Anime
	if err := db.Find(&animes).Error; err != nil {
		log.Printf("error retrieving anime from DB: %s", err)
		return
	}

	for _, anime := range animes {

		animeGOregexp, _ := regexp.Compile(`^https://animego.org/anime/.*$`)
		amediaOnline, _ := regexp.Compile(`^https://amedia.online/.*$`)

		if animeGOregexp.MatchString(anime.URL) {
			// If animego.org
			animegoorgcheck.Check(db, bot, anime)

		} else if amediaOnline.MatchString(anime.URL) {
			// If amedia.online
			amediaonlinecheck.Check(db, bot, anime)
		}

	}
}
