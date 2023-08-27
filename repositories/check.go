package repositories

import (
	amediaonlinecheck "anime-bot-schedule/checker/amedia.online"
	animegoorgcheck "anime-bot-schedule/checker/animego.org"
	animevostorgcheck "anime-bot-schedule/checker/animevost.org"
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
		animevostOrg, _ := regexp.Compile(`^https://animevost.org/tip/tv/.*$`)

		if animeGOregexp.MatchString(anime.URL) {
			animegoorgcheck.Check(db, bot, anime)
		} else if amediaOnline.MatchString(anime.URL) {
			amediaonlinecheck.Check(db, bot, anime)
		} else if animevostOrg.MatchString(anime.URL) {
			animevostorgcheck.Check(db, bot, anime)
		}
	}
}
