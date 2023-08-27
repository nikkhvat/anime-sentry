package parser

import (
	fouranimeisparsing "anime-bot-schedule/services/parser/4anime.is"
	"log"
)

func TestParser4animeIs() {
	data, err := fouranimeisparsing.Fetch("https://4anime.is/watch/masamunekuns-revenge-r-18419?ep=104316")

	if data != nil {
		log.Println("======================================")
		log.Println("- LastEpisode: (", data.LastEpisode, " )")
		log.Println("- TotalEpisodes: (", data.TotalEpisodes, " )")
		log.Println("- Poster: (", data.Poster, " )")
		log.Println("- Title: (", data.Title, " )")
		log.Println("- URL: (", data.URL, " )")
		log.Println("- LastEpisodeLink: (", data.LastEpisodeLink, " )")
		log.Println("======================================")
	}

	if err != nil {
		log.Println(err)
	}

}
