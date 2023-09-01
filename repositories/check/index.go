package repositories_check

import (
	repositories_animes "anime-bot-schedule/repositories/animes"

	fouranimeis_check "anime-bot-schedule/services/checker/4anime.is"
	amediaonline_check "anime-bot-schedule/services/checker/amedia.online"
	animegoorg_check "anime-bot-schedule/services/checker/animego.org"
	animevostorg_check "anime-bot-schedule/services/checker/animevost.org"

	fouranimeis_service "anime-bot-schedule/services/service/4anime.is"
	amediaonline_service "anime-bot-schedule/services/service/amedia.online"
	animegoorg_service "anime-bot-schedule/services/service/animego.org"
	animevostorg_service "anime-bot-schedule/services/service/animevost.org"

	"regexp"
)

func CheckAnimeStatus() {
	animes, err := repositories_animes.GetAnimes()

	if err != nil {
		return
	}

	for _, anime := range animes {

		animeGOregexp, _ := regexp.Compile(animegoorg_service.LINK_PATTERN)
		amediaOnline, _ := regexp.Compile(amediaonline_service.LINK_PATTERN)
		animevostOrg, _ := regexp.Compile(animevostorg_service.LINK_PATTERN)
		fouranimeIs, _ := regexp.Compile(fouranimeis_service.LINK_PATTERN)

		if animeGOregexp.MatchString(anime.URL) {
			animegoorg_check.Check(anime)
		} else if amediaOnline.MatchString(anime.URL) {
			amediaonline_check.Check(anime)
		} else if animevostOrg.MatchString(anime.URL) {
			animevostorg_check.Check(anime)
		} else if fouranimeIs.MatchString(anime.URL) {
			fouranimeis_check.Check(anime)
		}
	}
}
